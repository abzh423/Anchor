package api

import (
	"image"

	"github.com/anchormc/anchor/src/api/conf"
	"github.com/anchormc/anchor/src/api/game"
)

type Server interface {
	Initialize() error
	Start() error
	AcceptConnections()
	Favicon() image.Image
	GetAllClients() []Client
	PlayerCount() int
	RemoveClient(string)
	GetWorld(string) game.World
	CreateWorld(string, game.Generator) (game.World, error)
	GetConfig() *conf.Configuration
	Close() error
}
