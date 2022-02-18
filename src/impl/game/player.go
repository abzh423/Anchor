package game

import (
	"github.com/anchormc/anchor/src/api/game"
	"github.com/anchormc/protocol"
)

type Player struct {
	entityID int64
	username string
	uuid     string
	position protocol.AbsolutePosition
}

func NewPlayer(entityID int64, username, uuid string, position protocol.AbsolutePosition) *Player {
	return &Player{
		entityID: entityID,
		username: username,
		uuid:     uuid,
		position: position,
	}
}

func (p Player) EntityID() int64 {
	return p.entityID
}

func (p Player) Username() string {
	return p.username
}

func (p Player) UUID() string {
	return p.uuid
}

func (p Player) Position() protocol.AbsolutePosition {
	return p.position
}

var _ game.Player = &Player{}
