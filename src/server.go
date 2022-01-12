package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"path"

	"github.com/golangminecraft/minecraft-server/src/api/game/world"
	"github.com/golangminecraft/minecraft-server/src/api/server"
	"github.com/golangminecraft/minecraft-server/src/game"
	"github.com/golangminecraft/minecraft-server/src/game/generators"
	"github.com/golangminecraft/minecraft-server/src/networking"
	"github.com/golangminecraft/minecraft-server/src/protocol/handlers"
	"github.com/golangminecraft/minecraft-server/src/types"
)

type Server struct {
	id             string
	cwd            string
	isRunning      bool
	config         *types.Configuration
	clients        map[string]server.Client
	socket         server.Socket
	packetHandlers []server.PacketHandler
	privateKey     *rsa.PrivateKey
	worldManager   world.WorldManager
}

func NewServer(cwd string) *Server {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)

	if err != nil {
		log.Fatal(err)
	}

	id := make([]byte, 10)

	if _, err := rand.Read(id); err != nil {
		log.Fatal(err)
	}

	return &Server{
		id:             hex.EncodeToString(id),
		cwd:            cwd,
		isRunning:      false,
		config:         types.NewConfiguration(),
		clients:        make(map[string]server.Client),
		socket:         networking.NewSocket(),
		packetHandlers: make([]server.PacketHandler, 0),
		privateKey:     privateKey,
		worldManager:   game.NewWorldManager(),
	}
}

func (s *Server) Initialize() error {
	if err := s.config.ReadFromFile(path.Join(s.cwd, "config.yml")); err != nil {
		log.Println(err)
	}

	if err := s.config.WriteToFile(path.Join(s.cwd, "config.yml")); err != nil {
		log.Fatal(err)
	}

	if err := os.Mkdir(path.Join(s.cwd, "worlds"), 0666); err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}

	if err := os.Mkdir(path.Join(s.cwd, "logs"), 0666); err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}

	{
		generator := generators.NewFlatWorldGenerator(s.config.Seed)
		world := game.NewWorld("world", generator)
		s.worldManager.NewWorld(world)

		/* for x := int64(-3); x < 3; x++ {
			for z := int64(-3); z < 3; z++ {
				world.GenerateChunk(x, z)

				log.Printf("Generated chunk %d:%d for world %s\n", x, z, world.Name())
			}
		} */
	}

	s.packetHandlers = append(s.packetHandlers, handlers.HandshakeHandler{})
	s.packetHandlers = append(s.packetHandlers, handlers.PingHandler{})
	s.packetHandlers = append(s.packetHandlers, handlers.LoginStartHandler{})
	s.packetHandlers = append(s.packetHandlers, handlers.EncryptionResponseHandler{})
	s.packetHandlers = append(s.packetHandlers, handlers.ClientSettingsHandler{})
	s.packetHandlers = append(s.packetHandlers, handlers.RequestHandler{})

	return nil
}

func (s *Server) Start() error {
	if err := s.socket.Start(s.config.Host, s.config.Port); err != nil {
		return err
	}

	log.Printf("Server started on %s:%d\n", s.config.Host, s.config.Port)

	s.isRunning = true

	go (func() {
		for s.isRunning {
			client, err := s.socket.OnConnection()

			if err != nil {
				log.Println(err)

				continue
			}

			go (func() {
				if err := client.HandlePackets(s); err != nil {
					if !errors.Is(err, io.EOF) {
						log.Println(err)
					}

					delete(s.clients, client.ID())
				}
			})()

			s.clients[client.ID()] = client
		}
	})()

	return nil
}

func (s Server) GetPacketHandlers() []server.PacketHandler {
	return s.packetHandlers
}

func (s *Server) Stop() error {
	s.isRunning = false

	if err := s.socket.Stop(); err != nil {
		return err
	}

	log.Println("Socket server gracefully closed")

	return nil
}

func (s Server) GetCwd() string {
	return s.cwd
}

func (s Server) IsRunning() bool {
	return s.isRunning
}

func (s Server) GetClients() map[string]server.Client {
	return s.clients
}

func (s Server) GetSocket() server.Socket {
	return s.socket
}

func (s Server) OnlinePlayers() int {
	online := 0

	for _, v := range s.clients {
		if v.GetPlayer() != nil {
			online++
		}
	}

	return online
}

func (s Server) GetConfig() *types.Configuration {
	return s.config
}

func (s Server) GetPrivateKey() *rsa.PrivateKey {
	return s.privateKey
}

func (s Server) ID() string {
	return s.id
}

func (s Server) GetWorldManager() world.WorldManager {
	return s.worldManager
}

var _ server.Server = &Server{}
