package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	cwd string = ""
)

func init() {
	var err error

	if cwd, err = os.Getwd(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	server := NewServer(cwd)

	if err := server.Initialize(); err != nil {
		log.Fatal(err)
	}

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, os.Interrupt)

	go (func() {
		for server.Running() {
			var input string

			if _, err := fmt.Scanln(&input); err != nil {
				log.Fatal(err)
			}

			if err := server.ProcessConsoleCommand(input, &s); err != nil {
				log.Println(err)
			}
		}
	})()

	<-s

	if err := server.Close(); err != nil {
		log.Fatal(err)
	}

	log.Println("Server closed gracefully")
}
