package main

import (
	"errors"
	"fmt"
	"net"

	"github.com/anchormc/anchor/src/logger"
)

// NewSocket creates a new socket that contains the default values.
func NewSocket() *Socket {
	return &Socket{
		Listener: nil,
	}
}

// Socket is a utility struct for accepting incoming TCP connections from the
// network.
type Socket struct {
	Listener net.Listener
}

// Listen starts the socket server and allows incoming connections to the
// address specified by the arguments.
func (s *Socket) Listen(host string, port uint16) error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		return err
	}

	s.Listener = l

	return nil
}

// AcceptConnection blocks the thread until a new connection is received by the
// network.
func (s Socket) AcceptConnection() (net.Conn, error) {
	if s.Listener == nil {
		return nil, errors.New("socket: attempted to AcceptConnection() before starting the socket server")
	}

	return s.Listener.Accept()
}

// Close closes the socket connection and prevents any new incoming connections
// to the server.
func (s Socket) Close() error {
	if s.Listener == nil {
		logger.Warn("socket: Close() was called while the socket was already closed")

		return nil
	}

	return s.Listener.Close()
}
