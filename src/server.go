package main

import "github.com/anchormc/anchor/src/logger"

// NewServer creates a new server structure that contains the default values for
// starting a Minecraft server.
func NewServer() *Server {
	return &Server{
		Socket: NewSocket(),
	}
}

// Server is the root structure for interacting with the Minecraft server API,
// as well as performing all functions.
type Server struct {
	Socket *Socket
}

// Initialize prepares the server for startup, and does necessary tasks such as
// reading the configuration and creating the TCP socket.
func (s *Server) Initialize() (err error) {
	if err = logger.Initialize(); err != nil {
		return
	}

	return
}

// Start is called whenever the server is ready to accept members. This should
// always be called after Server#Initialize() is called.
func (s *Server) Start() error {
	return nil
}

// Close closes the server and all connections to the socket, and does any
// cleanup like saving chunk data and notifying any plugins of server close.
func (s *Server) Close() (err error) {
	if err = logger.Stop(); err != nil {
		return
	}

	return
}
