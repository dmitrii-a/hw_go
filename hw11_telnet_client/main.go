package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var t string
	flag.StringVar(&t, "timeout", "10s", "timeout")
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("host and port arguments not passed")
	}

	timeout, err := time.ParseDuration(t)
	if err != nil {
		log.Fatalf("Parsing time error")
	}

	address := net.JoinHostPort(flag.Args()[0], flag.Args()[1])
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatalf("connection failed: %s", err)
	}
	defer func(client TelnetClient) {
		err := client.Close()
		if err != nil {
			log.Fatalf("close connection failed: %s", err)
		}
	}(client)

	errCh := make(chan error)

	go func() {
		errCh <- client.Send()
	}()
	go func() {
		errCh <- client.Receive()
	}()

	quit := make(chan os.Signal, 1)
	defer close(quit)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT)

	select {
	case <-quit:
	case err := <-errCh:
		log.Printf("Error: %s", err)
	}
}
