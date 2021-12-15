package src

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"

	"github.com/golangminecraft/minecraft-server/src/logger"
	"gopkg.in/yaml.v2"
)

type Server struct {
	Logger   *logger.Logger
	Config   *Configuration
	Cwd      string
	Listener net.Listener
	Clients  []*Client
}

func NewServer() *Server {
	return &Server{
		Logger:   logger.NewLogger(1),
		Config:   &Configuration{},
		Cwd:      "",
		Listener: nil,
		Clients:  make([]*Client, 0),
	}
}

func (s *Server) Init() error {
	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	s.Cwd = cwd

	s.Config.LoadDefaults()

	data, err := ioutil.ReadFile(path.Join(cwd, "config.yml"))

	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	} else if err = yaml.Unmarshal(data, s.Config); err != nil {
		return err
	}

	data, err = yaml.Marshal(s.Config)

	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(path.Join(cwd, "config.yml"), data, 0660); err != nil {
		return err
	}

	return nil
}

func (s Server) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port))

	if err != nil {
		return err
	}

	defer l.Close()

	s.Listener = l

	s.Logger.Infof("Listening on %s:%d\n", s.Config.Host, s.Config.Port)

	go s.ReadStdin()

	for {
		conn, err := l.Accept()

		if err != nil {
			s.Logger.Error("Failed to accept connection:", err)

			continue
		}

		go s.HandleConnection(conn)
	}
}

func (s Server) ReadStdin() {

}

func (s Server) HandleConnection(conn net.Conn) {
	s.Logger.Infof("Received a connection from %s\n", conn.RemoteAddr())

	client := NewClient(conn)
	s.Clients = append(s.Clients, client)

	if err := client.Process(); err != nil {
		s.Logger.Errorf("Failed to process packet from client %s: %v\n", conn.RemoteAddr(), err)
	}

	conn.Close()
}

func (s Server) Close() error {
	for _, v := range s.Clients {
		if err := v.Close(); err != nil {
			return err
		}
	}

	return s.Listener.Close()
}
