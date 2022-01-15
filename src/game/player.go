package game

import (
	"math"

	"github.com/golangminecraft/minecraft-server/src/api"
	"github.com/golangminecraft/minecraft-server/src/api/enum"
	proto "github.com/golangminecraft/minecraft-server/src/api/protocol"
)

type Player struct {
	entityID int64
	gamemode enum.Gamemode
	uuid     proto.UUID
	username string
	position proto.AbsolutePosition
	rotation api.Rotation
}

func NewPlayer(entityID int64, username string, uuid proto.UUID, gamemode enum.Gamemode) api.Player {
	return &Player{
		entityID: entityID,
		gamemode: gamemode,
		username: username,
		uuid:     uuid,
		position: proto.AbsolutePosition{X: 0, Y: 63, Z: 0},
		rotation: api.Rotation{Yaw: 0, Pitch: math.Pi},
	}
}

func (p Player) EntityID() int64 {
	return p.entityID
}

func (p Player) Username() string {
	return p.username
}

func (p Player) UUID() proto.UUID {
	return p.uuid
}

func (p Player) Gamemode() enum.Gamemode {
	return p.gamemode
}

func (p Player) Position() proto.AbsolutePosition {
	return p.position
}

func (p Player) Rotation() api.Rotation {
	return p.rotation
}
