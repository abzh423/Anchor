package net

import (
	"errors"
	"fmt"
	"net"

	netAPI "github.com/anchormc/anchor/src/api/net"
)

type Socket struct {
	isRunning bool
	listener  net.Listener
}

func NewSocket() netAPI.Socket {
	return &Socket{
		listener: nil,
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

func (s *Socket) OnConnection() (net.Conn, error) {
	if !s.isRunning || s.listener == nil {
		return nil, errors.New("attempted to call OnConnection() when socket is not running")
	}

	return s.listener.Accept()
}

func (s *Socket) Close() error {
	if !s.isRunning || s.listener == nil {
		return errors.New("attempted to call Close() when socket is not running")
	}

	s.isRunning = false

	return s.listener.Close()
}

var _ netAPI.Socket = &Socket{}
