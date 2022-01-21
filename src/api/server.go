package api

import (
	"os"
	"time"

	"github.com/golangminecraft/minecraft-server/src/api/data"
	"github.com/golangminecraft/minecraft-server/src/api/enum"
	proto "github.com/golangminecraft/minecraft-server/src/api/protocol"
	"github.com/golangminecraft/minecraft-server/src/api/world"
)

type Server interface {
	Running() bool
	Initialize() error
	Start() error
	Close() error
	GetSocket() Socket
	AcceptConnections()
	ProcessConsoleCommand(string, *chan os.Signal) error
	AddClient(Client)
	RemoveClient(Client)
	OnlinePlayers() int
	MaxPlayers() int
	SamplePlayers() []data.StatusResponseSamplePlayer
	MOTD() proto.Chat
	Favicon() (*string, error)
	NextEntityID() int64
	OnlineMode() bool
	Difficulty() enum.Difficulty
	Hardcore() bool
	DefaultGamemode() enum.Gamemode
	DefaultWorld() *world.World
	ViewDistance() int
	SimulationDistance() int
	KeepAliveInterval() time.Duration
	Players() []Player
	Clients() []Client
	Host() string
	Port() uint16
	WorldManager() *world.WorldManager
}
