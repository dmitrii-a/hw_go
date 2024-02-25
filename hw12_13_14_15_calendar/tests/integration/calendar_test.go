package integration

import (
	"context"
	"testing"
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/application"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/domain"
	mq "github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/infrastructure/event"
	pb "github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/presentation/grpc/api/v1"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/tests"
	. "github.com/onsi/ginkgo" //nolint: revive
	. "github.com/onsi/gomega" //nolint: revive
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCart(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Calendar Suite")
}

var _ = Describe("Calendar API", func() {
	grpcEndpoint := common.GetServerAddr(
		common.Config.Server.GrpcHost,
		common.Config.Server.GrpcPort,
	)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(grpcEndpoint, opts...)
	Expect(err).ToNot(HaveOccurred())
	grpcClient := pb.NewEventServiceV1Client(conn)
	var (
		ctx   context.Context
		event *domain.Event
	)
	BeforeEach(func() {
		ctx = context.Background()
		event = tests.GenerateTestEvent()
	})
	Describe("Create Event", func() {
		It("creating an event", func() {
			e, err := grpcClient.CreateEvent(ctx, tests.CreateTestEventRequest(event))
			Expect(err).ToNot(HaveOccurred())
			Expect(e).ToNot(BeNil())
			Expect(e.Event).ToNot(BeNil())
			Expect(e.Event.Id).ToNot(BeEmpty())
			Expect(e.Event.Title).To(Equal(event.Title))
			Expect(e.Event.StartTime.AsTime()).To(Equal(event.StartTime))
			Expect(e.Event.EndTime.AsTime()).To(Equal(*event.EndTime))
			Expect(e.Event.NotifyTime.AsTime()).To(Equal(*event.NotifyTime))
			Expect(e.Event.UserId).To(Equal(event.UserID))
		})
		It("error creating an event(end time less start time)", func() {
			event.StartTime = event.EndTime.Add(time.Hour)
			_, err := grpcClient.CreateEvent(ctx, tests.CreateTestEventRequest(event))
			Expect(err).Should(HaveOccurred())
			Expect(
				err.Error(),
			).To(Equal("rpc error: code = InvalidArgument desc = end time must be greater than start time"))
		})
		It("error creating an event(notify time less start time)", func() {
			t := event.StartTime.Add(-time.Hour)
			event.NotifyTime = &t
			_, err := grpcClient.CreateEvent(ctx, tests.CreateTestEventRequest(event))
			Expect(err).Should(HaveOccurred())
			Expect(
				err.Error(),
			).To(Equal("rpc error: code = InvalidArgument desc = notify time must be greater than start time"))
		})
		It("creating an empty event", func() {
			event := &domain.Event{}
			_, err := grpcClient.CreateEvent(ctx, tests.CreateTestEventRequest(event))
			Expect(err).Should(HaveOccurred())
		})
	})
	Describe("Update Event", func() {
		It("updating an event", func() {
			e, err := grpcClient.CreateEvent(ctx, tests.CreateTestEventRequest(event))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(e).ToNot(BeNil())
			event.Title = "test"
			event.ID = e.Event.Id
			createdTime := e.Event.CreatedTime.AsTime()
			event.CreatedTime = &createdTime
			e, err = grpcClient.UpdateEvent(ctx, tests.CreateTestEventRequest(event))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(e).ToNot(BeNil())
			Expect(e.Event).ToNot(BeNil())
			Expect(e.Event.Id).ToNot(BeEmpty())
			Expect(e.Event.Title).To(Equal(event.Title))
			Expect(e.Event.StartTime.AsTime()).To(Equal(event.StartTime))
			Expect(e.Event.EndTime.AsTime()).To(Equal(*event.EndTime))
			Expect(e.Event.NotifyTime.AsTime()).To(Equal(*event.NotifyTime))
			Expect(e.Event.UserId).To(Equal(event.UserID))
			Expect(e.Event.CreatedTime.AsTime()).To(Equal(*event.CreatedTime))
		})
		It("updating an event that doesn't exist", func() {
			_, err := grpcClient.UpdateEvent(ctx, tests.CreateTestEventRequest(event))
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("rpc error: code = NotFound desc = event doesn't exist"))
		})
	})
	Describe("Delete Event", func() {
		It("deleting an event", func() {
			e, err := grpcClient.CreateEvent(ctx, tests.CreateTestEventRequest(event))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(e).ToNot(BeNil())
			_, err = grpcClient.DeleteEvent(ctx, &pb.EventIDRequest{Id: e.Event.Id})
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("deleting an event with incorrect id", func() {
			_, err := grpcClient.DeleteEvent(ctx, &pb.EventIDRequest{Id: "wrong"})
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(
				Equal(
					"rpc error: code = Unknown desc = invalid EventIDRequest.Id: " +
						"value must be a valid UUID | caused by: invalid uuid format",
				),
			)
		})
		It("deleting a non-existent element", func() {
			_, err := grpcClient.DeleteEvent(ctx, &pb.EventIDRequest{Id: event.ID})
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("rpc error: code = NotFound desc = event doesn't exist"))
		})
	})
	Describe("Range events", func() {
		It("getting events by range", func() {
			_, err := grpcClient.CreateEvent(ctx, tests.CreateTestEventRequest(event))
			Expect(err).ShouldNot(HaveOccurred())
			events, err := grpcClient.GetEventsByPeriod(ctx, &pb.TimePeriodRequest{
				StartTime: timestamppb.New(event.StartTime),
				EndTime:   timestamppb.New(*event.EndTime),
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(events).ToNot(BeNil())
			Expect(events.Events).ToNot(BeNil())
			Expect(len(events.Events)).To(Equal(1))
		})
	})
	Describe("Test notification", func() {
		It("testing notification", func() {
			t := time.Now()
			event.NotifyTime = &t
			e, err := grpcClient.CreateEvent(ctx, tests.CreateTestEventRequest(event))
			Expect(err).ShouldNot(HaveOccurred())
			time.Sleep(time.Duration(common.Config.Scheduler.PublishPeriodTime) * time.Second)
			client := mq.NewRabbitClient()
			ch, err := client.Consume(application.EventResultQueueName)
			Expect(err).ShouldNot(HaveOccurred())
			msg, ok := <-ch
			Expect(ok).To(BeTrue())
			Expect(msg).ToNot(BeNil())
			Expect(msg).To(Equal([]byte(e.Event.Id)))
		})
	})
})
