package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/golangminecraft/minecraft-server/src/api"
	"github.com/golangminecraft/minecraft-server/src/api/data"
	"github.com/golangminecraft/minecraft-server/src/api/enum"
	proto "github.com/golangminecraft/minecraft-server/src/api/protocol"
	"github.com/golangminecraft/minecraft-server/src/api/world"
	"github.com/golangminecraft/minecraft-server/src/components"
	"github.com/golangminecraft/minecraft-server/src/game/generators"
	"github.com/golangminecraft/minecraft-server/src/game/store"
	"github.com/golangminecraft/minecraft-server/src/handlers"
	"github.com/golangminecraft/minecraft-server/src/net"
	"github.com/golangminecraft/minecraft-server/src/query"
)

type Server struct {
	isRunning   bool
	config      api.Configuration
	cwd         string
	socket      api.Socket
	clients     map[string]api.Client
	worlds      map[string]*world.World
	favicon     *image.Image
	entityID    int64
	queryServer api.QueryServer
}

func NewServer(cwd string) api.Server {
	return &Server{
		isRunning:   false,
		config:      api.NewConfiguration(),
		cwd:         cwd,
		socket:      net.NewSocket(),
		clients:     make(map[string]api.Client),
		worlds:      make(map[string]*world.World),
		favicon:     nil,
		entityID:    0,
		queryServer: query.NewServer(),
	}
}

func (s *Server) Initialize() error {
	if err := s.config.ReadFromFile(path.Join(s.cwd, "config.yml")); err != nil {
		log.Println(err)
	}

	if err := s.config.WriteToFile(path.Join(s.cwd, "config.yml")); err != nil {
		log.Fatal(err)
	}

	if err := s.config.Validate(); err != nil {
		log.Fatal(err)
	}

	if err := s.queryServer.Initialize(s); err != nil {
		return err
	}

	if err := os.Mkdir(path.Join(s.cwd, "worlds"), 0777); err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}

	if err := os.Mkdir(path.Join(s.cwd, "logs"), 0777); err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}

	if data, err := ioutil.ReadFile(path.Join(s.cwd, "favicon.png")); err == nil {
		img, err := png.Decode(bytes.NewReader(data))

		if err != nil {
			return err
		}

		s.favicon = &img
	}

	for _, component := range components.Components {
		if err := component.Initialize(s); err != nil {
			log.Fatal(err)
		}
	}

	for _, worldMeta := range s.config.Worlds {
		var worldGenerator world.WorldGenerator
		var worldStore world.WorldStore

		switch worldMeta.Generator {
		case "flat":
			{
				worldGenerator = generators.FlatGenerator

				break
			}
		default:
			{
				return fmt.Errorf("unknown generator \"%s\" for world \"%s\"", worldMeta.Generator, worldMeta.Name)
			}
		}

		switch worldMeta.Store.Name {
		case "sqlite":
			{
				worldStore = store.NewSQLiteStore()

				break
			}
		default:
			{
				return fmt.Errorf("unknown store \"%s\" for world \"%s\"", worldMeta.Store.Name, worldMeta.Name)
			}
		}

		world, err := s.NewWorld(worldMeta.Name, worldGenerator, worldStore, worldMeta.Store.Options)

		if err != nil {
			return err
		}

		for x := -10; x <= 10; x++ {
			for z := -10; z <= 10; z++ {
				if err := world.GenerateChunk(int64(x), int64(z)); err != nil {
					return err
				}
			}
		}

		log.Printf("Created new world %s\n", worldMeta.Name)
	}

	log.Printf("Initialized %d components\n", len(components.Components))
	log.Printf("Loaded %d packet handlers\n", len(handlers.Handlers))

	return nil
}

func (s *Server) Start() error {
	if err := s.socket.Start(s.config.Host, s.config.Port); err != nil {
		return err
	}

	for _, component := range components.Components {
		if err := component.Start(); err != nil {
			log.Fatal(err)
		}
	}

	s.isRunning = true

	if s.config.EnableQuery {
		if err := s.queryServer.Start(s.config.QueryHost, s.config.QueryPort); err != nil {
			return err
		}

		log.Printf("Query server running on %s:%d\n", s.config.QueryHost, s.config.QueryPort)
	}

	log.Printf("Server running on %s:%d\n", s.config.Host, s.config.Port)

	go s.AcceptConnections()

	return nil
}

func (s *Server) Close() error {
	s.isRunning = false

	for _, client := range s.clients {
		// TODO properly send shutdown message to players

		if err := client.Close(); err != nil {
			log.Println(err)
		}
	}

	for _, world := range s.worlds {
		log.Printf("Saving all chunks for world %s...\n", world.Name())

		if err := world.Close(); err != nil {
			return err
		}

		log.Printf("Saved all chunks for world %s\n", world.Name())
	}

	if s.config.EnableQuery {
		if err := s.queryServer.Close(); err != nil {
			return err
		}

		log.Println("Closed query server")
	}

	return s.socket.Close()
}

func (s *Server) AcceptConnections() {
	for s.isRunning {
		client, err := s.socket.AcceptConnection()

		if err != nil {
			log.Println(err)

			continue
		}

		log.Printf("Received a connection from %s\n", client.RemoteAddr())

		s.AddClient(client)

		go (func() {
			if err := client.HandleConnection(s); err != nil {
				log.Println(err)

				if err = client.Close(); err != nil {
					log.Println(err)
				}
			}

			s.RemoveClient(client)
		})()
	}
}

func (s *Server) AddClient(client api.Client) {
	s.clients[client.ID()] = client
}

func (s *Server) RemoveClient(client api.Client) {
	delete(s.clients, client.ID())
}

func (s Server) GetSocket() api.Socket {
	return s.socket
}

func (s Server) OnlinePlayers() int {
	onlinePlayers := 0

	for _, client := range s.clients {
		if client.GetPlayer() == nil {
			continue
		}

		onlinePlayers++
	}

	return onlinePlayers
}

func (s Server) MaxPlayers() int {
	return s.config.MaxPlayers
}

func (s Server) SamplePlayers() []data.StatusResponseSamplePlayer {
	return make([]data.StatusResponseSamplePlayer, 0)
}

func (s Server) Favicon() (*string, error) {
	if s.favicon == nil {
		return nil, nil
	}

	buf := &bytes.Buffer{}

	if err := png.Encode(buf, *s.favicon); err != nil {
		return nil, err
	}

	value := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	return &value, nil
}

func (s Server) MOTD() proto.Chat {
	return proto.Chat{
		Text: s.config.MOTD,
	}
}

func (s Server) Running() bool {
	return s.isRunning
}

func (s Server) ProcessConsoleCommand(command string, shutdown *chan os.Signal) error {
	args := strings.Split(command, " ")

	switch args[0] {
	case "stop", "shutdown", "close":
		{
			*shutdown <- os.Interrupt

			return nil
		}
	default:
		{
			return fmt.Errorf("unknown console command: %s", args[0])
		}
	}
}

func (s *Server) NextEntityID() int64 {
	s.entityID++

	return s.entityID
}

func (s Server) OnlineMode() bool {
	return s.config.OnlineMode
}

func (s Server) Difficulty() enum.Difficulty {
	switch s.config.Difficulty {
	case "peaceful":
		{
			return enum.DifficultyPeaceful
		}
	case "easy":
		{
			return enum.DifficultyEasy
		}
	case "hard":
		{
			return enum.DifficultyHard
		}
	default:
		{
			return enum.DifficultyNormal
		}
	}
}

func (s Server) Hardcore() bool {
	return s.config.Hardcore
}

func (s Server) DefaultGamemode() enum.Gamemode {
	switch s.config.DefaultGamemode {
	case "creative":
		{
			return enum.GamemodeCreative
		}
	case "adventure":
		{
			return enum.GamemodeAdventure
		}
	case "spectator":
		{
			return enum.GamemodeSpectator
		}
	default:
		{
			return enum.GamemodeSurvival
		}
	}
}

func (s Server) WorldCount() int {
	return len(s.worlds)
}

func (s Server) GetWorld(name string) *world.World {
	return s.worlds[name]
}

func (s *Server) NewWorld(name string, generator world.WorldGenerator, store world.WorldStore, storeOptions map[string]interface{}) (*world.World, error) {
	if _, ok := s.worlds[name]; ok {
		return nil, fmt.Errorf("unable to create a new world, already have one loaded with name %s", name)
	}

	world := world.NewWorld(name, path.Join(s.cwd, "worlds"), store, generator)

	if err := world.Initialize(storeOptions); err != nil {
		return nil, err
	}

	s.worlds[name] = world

	return world, nil
}

func (s Server) ViewDistance() int {
	return s.config.ViewDistance
}

func (s Server) SimulationDistance() int {
	return s.config.SimulationDistance
}

func (s Server) KeepAliveInterval() time.Duration {
	return time.Duration(s.config.KeepAliveInterval) * time.Second
}

func (s Server) Players() []api.Player {
	players := make([]api.Player, 0)

	for _, client := range s.clients {
		player := client.GetPlayer()

		if player != nil {
			players = append(players, player)
		}
	}

	return players
}

func (s Server) Clients() []api.Client {
	clients := make([]api.Client, 0)

	for _, client := range s.clients {
		clients = append(clients, client)
	}

	return clients
}

func (s Server) Host() string {
	return s.config.Host
}

func (s Server) Port() uint16 {
	return s.config.Port
}

var _ api.Server = &Server{}
