package net

import (
	"fmt"
	"net"

	"github.com/golangminecraft/minecraft-server/src/api"
)

type Socket struct {
	Listener net.Listener
}

func NewSocket() api.Socket {
	return &Socket{
		Listener: nil,
	}
}

func (s *Socket) Start(host string, port uint16) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		return err
	}

	s.Listener = listener

	return nil
}

func (s *Socket) Close() error {
	return s.Listener.Close()
}

func (s *Socket) AcceptConnection() (api.Client, error) {
	conn, err := s.Listener.Accept()

	if err != nil {
		return nil, err
	}

	client, err := NewClient(conn)

	if err != nil {
		if err := conn.Close(); err != nil {
			return nil, err
		}

		return nil, err
	}

	return client, nil
}

var _ api.Socket = &Socket{}
