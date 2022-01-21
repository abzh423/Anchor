package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"

	log "github.com/golangminecraft/minecraft-server/src/api/logger"
)

var (
	cwd            string      = ""
	previousLogDay int         = -1
	logFile        *os.File    = nil
	logMutex       *sync.Mutex = &sync.Mutex{}
)

func init() {
	var err error

	if cwd, err = os.Getwd(); err != nil {
		panic(err)
	}
}

func main() {
	server, err := NewServer(cwd)

	if err != nil {
		panic(err)
	}

	if err := server.Initialize(); err != nil {
		panic(err)
	}

	if err := server.Start(); err != nil {
		panic(err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, os.Interrupt)

	go (func() {
		if err := recover(); err != nil {
			log.Errorf("main", "recovered fatal error in main thread: %v\n", err)

			s <- os.Interrupt
		}

		// TODO recover panics in goroutines
	})()

	go (func() {
		for server.Running() {
			var input string

			if _, err := fmt.Scanln(&input); err != nil {
				panic(err)
			}

			if err := server.ProcessConsoleCommand(input, &s); err != nil {
				log.Error("main", err)
			}
		}
	})()

	go (func() {
		if err := os.Mkdir(path.Join(cwd, "logs"), 0777); err != nil && !errors.Is(err, os.ErrExist) {
			panic(err)
		}

		for {
			message := <-*log.OnMessage

			logMutex.Lock()

			if time.Now().Day() != previousLogDay {
				if logFile != nil {
					if err := logFile.Close(); err != nil {
						panic(err)
					}
				}

				f, err := os.OpenFile(path.Join(cwd, "logs", time.Now().Format("01-02-2006.log")), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)

				if err != nil {
					panic(err)
				}

				logFile = f
			}

			if _, err := logFile.WriteString(message); err != nil {
				panic(err)
			}

			logMutex.Unlock()
		}
	})()

	<-s

	if err := server.Close(); err != nil {
		panic(err)
	}

	log.Info("main", "Server closed gracefully")
}
