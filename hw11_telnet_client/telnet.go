package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address    string
	timeout    time.Duration
	connection net.Conn
	in         io.ReadCloser
	out        io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *telnetClient) Connect() error {
	if c.in == nil {
		return fmt.Errorf("in is invalid")
	}
	if c.out == nil {
		return fmt.Errorf("out is invalid")
	}
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.connection = conn
	return nil
}

func read(reader io.Reader) ([]byte, error) {
	data := make([]byte, 1024)
	n, err := reader.Read(data)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, err
		}
		return nil, fmt.Errorf("receiving error: %w", err)
	}
	return data[:n], nil
}

func (c *telnetClient) Close() error {
	return c.connection.Close()
}

func (c *telnetClient) Send() error {
	reader := bufio.NewReader(c.in)
	for {
		data, err := read(reader)
		if data == nil {
			return nil
		}
		if err != nil {
			return err
		}
		if _, err := c.connection.Write(data); err != nil {
			return fmt.Errorf("sending error: %w", err)
		}
	}
}

func (c *telnetClient) Receive() error {
	reader := bufio.NewReader(c.connection)
	for {
		data, err := read(reader)
		if data == nil {
			return nil
		}
		if err != nil {
			return err
		}
		if _, err := c.out.Write(data); err != nil {
			return fmt.Errorf("receiving error: %w", err)
		}
	}
}
