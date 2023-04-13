package main

import (
	"os"
	"os/signal"
)

func main() {
	server := NewServer()

	if err := server.Initialize(); err != nil {
		panic(err)
	}

	if err := server.Start(); err != nil {
		panic(err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	<-s

	if err := server.Close(); err != nil {
		panic(err)
	}
}
