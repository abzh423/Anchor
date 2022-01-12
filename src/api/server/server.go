package server

import (
	"crypto/rsa"

	"github.com/golangminecraft/minecraft-server/src/api/game/world"
	"github.com/golangminecraft/minecraft-server/src/types"
)

type Server interface {
	Initialize() error
	ID() string
	Start() error
	Stop() error
	IsRunning() bool
	GetCwd() string
	GetClients() map[string]Client
	GetSocket() Socket
	GetPacketHandlers() []PacketHandler
	GetConfig() *types.Configuration
	OnlinePlayers() int
	GetPrivateKey() *rsa.PrivateKey
	GetWorldManager() world.WorldManager
}
