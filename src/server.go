package src

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"

	"github.com/golangminecraft/minecraft-server/src/logger"
	"gopkg.in/yaml.v2"
)

type Server struct {
	IsRunning bool
	Logger    *logger.Logger
	Config    *Configuration
	Cwd       string
	Listener  net.Listener
	Clients   map[string]*Client
}

func NewServer() *Server {
	return &Server{
		IsRunning: false,
		Logger:    logger.NewLogger(1),
		Config:    &Configuration{},
		Cwd:       "",
		Listener:  nil,
		Clients:   make(map[string]*Client),
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

func (s *Server) Start() error {
	s.IsRunning = true

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port))

	if err != nil {
		return err
	}

	defer l.Close()

	s.Listener = l

	s.Logger.Infof("Listening on %s:%d\n", s.Config.Host, s.Config.Port)

	go s.ReadStdin()

	for s.IsRunning {
		conn, err := l.Accept()

		if err != nil {
			s.Logger.Error("Failed to accept connection:", err)

			continue
		}

		go s.HandleConnection(conn)
	}

	return nil
}

func (s *Server) ReadStdin() {
	buf := bufio.NewReader(os.Stdin)

	for s.IsRunning {
		data, err := buf.ReadBytes(0x0A)

		if err != nil {
			panic(err)
		}

		args := strings.Split(strings.Trim(string(data[:len(data)-1]), " "), " ")

		if len(args) < 1 {
			continue
		}

		switch args[0] {
		case "stop":
			{
				s.Logger.Info("Stopping the server...")

				s.IsRunning = false

				if err = s.Close(); err != nil {
					panic(err)
				}

				break
			}
		}
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	s.Logger.Infof("Received a connection from %s\n", conn.RemoteAddr())

	client := NewClient(s, conn)

	s.Clients[client.UUID] = client

	if err := client.Process(); err != nil {
		if !errors.Is(err, io.EOF) {
			s.Logger.Errorf("Failed to process packet from client %s: %v\n", conn.RemoteAddr(), err)
		}
	}

	delete(s.Clients, client.UUID)

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
