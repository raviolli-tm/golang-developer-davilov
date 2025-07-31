package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "timeout")
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 2 {
		log.Fatal("Usage: go run main.go --timeout=<timeout> <host> <port>")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGINT)
	defer cancel()

	client := NewTelnetClient(fmt.Sprintf("%s:%s", host, port), timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		log.Fatalf("Unable to connect by tcp, %v", err)
	}

	go func() {
		err := client.Receive()
		if err != nil {
			log.Printf("Unable to receive from tcp, %v \n", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			err := client.Close()
			if err != nil {
				log.Fatalf("Unable to close tcp connection, %v", err)
			}
		default:
			err = client.Send()
			if err != nil {
				cancel()
			}
		}

	}

}
