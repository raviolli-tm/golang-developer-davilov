package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

type result struct {
	operation string
	err       error
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go-telnet [--timeout=10s] <host> <port>")
		os.Exit(2)
	}

	if err := run(); err != nil {
		log.Printf("go-telnet stopped with error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	host, port := flag.Arg(0), flag.Arg(1)
	address := net.JoinHostPort(host, port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("unable to close connection: %v", err)
		}
	}()

	log.Printf("...Connected to %s", address)

	results := make(chan result, 2)

	go func() {
		results <- result{operation: "send", err: client.Send()}
	}()
	go func() {
		results <- result{operation: "receive", err: client.Receive()}
	}()

	select {
	case <-ctx.Done():
		log.Println("...Interrupted")
	case res := <-results:
		switch {
		case res.err != nil:
			log.Printf("...%s error: %v", res.operation, res.err)
		case res.operation == "send":
			log.Println("...EOF")
		default:
			log.Println("...Connection was closed by peer")
		}
	}

	return nil
}
