package api

import (
	"os"
	"time"

	"github.com/anchormc/anchor/src/api/data"
	"github.com/anchormc/anchor/src/api/enum"
	proto "github.com/anchormc/anchor/src/api/protocol"
	"github.com/anchormc/anchor/src/api/world"
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
	DefaultWorld() world.World
	ViewDistance() int
	SimulationDistance() int
	KeepAliveInterval() time.Duration
	Players() []Player
	Clients() []Client
	Host() string
	Port() uint16
	WorldManager() world.WorldManager
}
