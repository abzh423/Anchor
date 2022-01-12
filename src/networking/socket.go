package networking

import (
	"fmt"
	"log"
	"net"

	"github.com/golangminecraft/minecraft-server/src/api/server"
)

type Socket struct {
	isRunning bool
	listener  net.Listener
}

func NewSocket() *Socket {
	return &Socket{
		isRunning: false,
		listener:  nil,
	}
}

func (s Socket) IsRunning() bool {
	return s.isRunning
}

func (s *Socket) Start(host string, port uint16) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		return err
	}

	s.listener = listener
	s.isRunning = true

	return nil
}

func (s *Socket) Stop() error {
	s.isRunning = false

	return s.listener.Close()
}

func (s Socket) OnConnection() (server.Client, error) {
	conn, err := s.listener.Accept()

	if err != nil {
		log.Printf("Failed to accept connection: %s\n", err)

		return nil, err
	}

	log.Printf("Received a connection from %s\n", conn.RemoteAddr())

	return NewClient(conn), nil
}

var _ server.Socket = &Socket{}
