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
	// waiting for interrupt signal (Ctrl + C)
	signal.Notify(s, os.Interrupt)
	// blocks execution until we receive the signal
	<-s

	if err := server.Close(); err != nil {
		panic(err)
	}
}
