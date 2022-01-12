package game

import (
	"bytes"
	"encoding/binary"
	"strconv"

	"github.com/golangminecraft/minecraft-server/src/api/enum"
	"github.com/golangminecraft/minecraft-server/src/api/game"
	"github.com/golangminecraft/minecraft-server/src/api/game/world"
	"github.com/golangminecraft/minecraft-server/src/api/protocol"
)

type Player struct {
	entityID int
	uuid     string
	username string
	position protocol.AbsolutePosition
	rotation game.Rotation
}

func NewPlayer(entityID int, username string) *Player {
	return &Player{
		entityID: entityID,
		uuid:     "",
		username: username,
		position: protocol.AbsolutePosition{
			X: 0,
			Y: 63,
			Z: 0,
		},
	}
}

func (p Player) Type() enum.EntityType {
	return enum.EntityTypePlayer
}

func (p *Player) SetUUID(uuid string) {
	p.uuid = uuid
}

func (p Player) UUID() string {
	return p.uuid
}

func (p Player) EntityID() int {
	return p.entityID
}

func (p Player) UUIDBytes() ([]byte, error) {
	firstHalfParsed, err := strconv.ParseUint(p.uuid[0:16], 16, 64)

	if err != nil {
		return nil, err
	}

	secondHalfParsed, err := strconv.ParseUint(p.uuid[16:32], 16, 64)

	if err != nil {
		return nil, err
	}

	buffer := &bytes.Buffer{}

	if err := binary.Write(buffer, binary.BigEndian, firstHalfParsed); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, secondHalfParsed); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (p Player) UUIDWithHyphens() string {
	return p.uuid
}

func (p Player) Username() string {
	return p.username
}

func (p Player) Position() protocol.AbsolutePosition {
	return p.position
}

func (p Player) Rotation() game.Rotation {
	return p.rotation
}

var _ game.Player = &Player{}
var _ world.Entity = &Player{}
