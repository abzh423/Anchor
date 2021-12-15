package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/golangminecraft/minecraft-server/src"
)

var (
	server *src.Server = src.NewServer()
)

func init() {
	if err := server.Init(); err != nil {
		panic(err)
	}
}

func main() {
	go (func() {
		if err := recover(); err != nil {
			server.Logger.Errorf("Recovered error from server: %v\n", err)

			if err = server.Close(); err != nil {
				server.Logger.Fatalf("Failed to close server after recovering: %v\n", err)
			}

			os.Exit(1)
		}
	})()

	go (func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		server.Logger.Info("Received SIGTERM, closing the server...")

		if err := server.Close(); err != nil {
			server.Logger.Errorf("Failed to close server after receiving SIGTERM: %v\n", err)
		}

		os.Exit(0)
	})()

	panic(server.Start())
}
