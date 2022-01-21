package game

import (
	"fmt"
	"math"

	"github.com/golangminecraft/minecraft-server/src/api"
	"github.com/golangminecraft/minecraft-server/src/api/enum"
	proto "github.com/golangminecraft/minecraft-server/src/api/protocol"
)

type Player struct {
	entityID   int64
	gamemode   enum.Gamemode
	uuid       proto.UUID
	username   string
	position   proto.AbsolutePosition
	rotation   api.Rotation
	sentChunks map[string]bool
}

func NewPlayer(entityID int64, username string, uuid proto.UUID, gamemode enum.Gamemode) api.Player {
	return &Player{
		entityID:   entityID,
		gamemode:   gamemode,
		username:   username,
		uuid:       uuid,
		position:   proto.AbsolutePosition{X: 0, Y: 63, Z: 0},
		rotation:   api.Rotation{Yaw: 0, Pitch: math.Pi},
		sentChunks: make(map[string]bool),
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

func (p *Player) SetPosition(pos proto.AbsolutePosition) {
	p.position = pos
}

func (p Player) Chunk() proto.RelativePosition {
	return proto.RelativePosition{
		X: int64(p.position.X) / 16,
		Y: 0,
		Z: int64(p.position.Z) / 16,
	}
}

func (p *Player) SetSentChunk(x, z int64) {
	p.sentChunks[fmt.Sprintf("%d:%d", x, z)] = true
}

func (p Player) HasSentChunk(x, z int64) bool {
	_, ok := p.sentChunks[fmt.Sprintf("%d:%d", x, z)]

	return ok
}

func (p Player) Rotation() api.Rotation {
	return p.rotation
}

func (p *Player) SetRotation(rot api.Rotation) {
	p.rotation = rot
}
