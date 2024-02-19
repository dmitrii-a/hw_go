package event

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/domain"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitClient struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	initConnCh chan interface{}
	done       chan interface{}
	active     bool
}

// RabbitClient is the interface for the rabbitmq client.
type RabbitClient interface {
	domain.EventProducer
	domain.EventConsumer
}

// NewRabbitClient returns a new instance of the rabbitmq client.
func NewRabbitClient() RabbitClient {
	client := &rabbitClient{}
	client.done = make(chan interface{})
	client.initConnCh = make(chan interface{}, 1)
	client.initClient()
	return client
}

func (client *rabbitClient) setConnect() error {
	common.Logger.Info().Msg("starting the connection setup")
	var c string
	if common.Config.RabbitMQ.Username == "" {
		c = fmt.Sprintf("amqp://%s:%v/", common.Config.RabbitMQ.Host, common.Config.RabbitMQ.Port)
	} else {
		c = fmt.Sprintf(
			"amqp://%s:%s@%s:%v/",
			common.Config.RabbitMQ.Username,
			common.Config.RabbitMQ.Password,
			common.Config.RabbitMQ.Host,
			common.Config.RabbitMQ.Port,
		)
	}
	conn, err := amqp.Dial(c)
	if common.IsErr(err) {
		return err
	}
	client.conn = conn
	return err
}

func (client *rabbitClient) setChannel() error {
	common.Logger.Info().Msg("starting the channel setup")
	var err error
	client.channel, err = client.conn.Channel()
	return err
}

func (client *rabbitClient) listenNotify() {
	common.Logger.Info().Msg("start listening rabbitmq notification")
	connClose := client.conn.NotifyClose(make(chan *amqp.Error))
	connBlocked := client.conn.NotifyBlocked(make(chan amqp.Blocking))
	chClose := client.channel.NotifyClose(make(chan *amqp.Error))
	var err error
	go func() {
		for {
			select {
			case <-client.done:
				common.Logger.Info().Msg("stop listening rabbitmq notification")
				return
			case <-connBlocked:
				common.Logger.Error().Msg("connection blocked")
				client.initClient()
			case err = <-connClose:
				common.Logger.Error().Msgf("connection closed: %v", err)
				client.initClient()
			case err = <-chClose:
				common.Logger.Error().Msgf("channel closed: %v", err)
				client.initClient()
			}
		}
	}()
}

func (client *rabbitClient) initClient() {
initLoop:
	for {
		select {
		case <-client.done:
			return
		default:
			client.active = false
			err := client.setConnect()
			if common.IsErr(err) {
				common.Logger.Error().Msgf("rabbitmq conntection error: %v", err)
				time.Sleep(time.Second)
				continue
			}
			err = client.setChannel()
			if common.IsErr(err) {
				common.Logger.Error().Msgf("rabbitmq create channel error: %v", err)
				time.Sleep(time.Second)
				continue
			}
			client.active = true
			break initLoop
		}
	}
	client.listenNotify()
	client.initConnCh <- struct{}{}
	common.Logger.Info().Msg("rabbitmq client is active")
}

// Consume consumes messages from the queue.
func (client *rabbitClient) Consume(name string) (<-chan []byte, error) {
	ch := make(chan []byte)
	var msg <-chan amqp.Delivery
	go func() {
		for {
			if !client.active {
				continue
			}
			select {
			case <-client.done:
				return
			case <-client.initConnCh:
				for {
					q, err := client.channel.QueueDeclare(
						name,
						true,
						false,
						false,
						false,
						amqp.Table{"x-queue-mode": "lazy"},
					)
					if common.IsErr(err) {
						common.Logger.Error().Msgf("failed to declare a queue: %v", err)
					}
					msg, err = client.channel.Consume(
						q.Name,
						uuid.New().String(),
						false,
						false,
						false,
						false,
						nil,
					)
					if common.IsErr(err) {
						common.Logger.Error().Msgf("failed to consume: %v", err)
						continue
					}
					break
				}
			case d := <-msg:
				common.Logger.Info().Msgf("received a message: %s", d.Body)
				ch <- d.Body
				err := d.Ack(false)
				if common.IsErr(err) {
					common.Logger.Error().Msgf("failed to ack: %v", err)
				}
			}
		}
	}()
	return ch, nil
}

// Publish publishes a message to the queue.
func (client *rabbitClient) Publish(ctx context.Context, queueName string, data []byte) error {
	if !client.active {
		return fmt.Errorf("rabbitmq client is not active")
	}
	q, err := client.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		amqp.Table{"x-queue-mode": "lazy"},
	)
	if common.IsErr(err) {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}
	err = client.channel.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			MessageId:    uuid.New().String(),
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         data,
		})
	if common.IsErr(err) {
		return fmt.Errorf("failed to publish to queue: %w", err)
	}
	common.Logger.Info().Msgf("published msg to queue \"%s\"", queueName)
	return nil
}

// Close closes the rabbitmq client.
func (client *rabbitClient) Close() error {
	err := client.channel.Close()
	if common.IsErr(err) {
		return err
	}
	err = client.conn.Close()
	if common.IsErr(err) {
		return err
	}
	close(client.done)
	return nil
}
