package main

import (
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
	<-s

	if err := server.Stop(); err != nil {
		log.Fatal(err)
	}
}
