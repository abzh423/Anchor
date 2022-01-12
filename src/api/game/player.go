package game

import "github.com/golangminecraft/minecraft-server/src/api/protocol"

type Player interface {
	Username() string
	EntityID() int
	UUID() string
	UUIDWithHyphens() string
	UUIDBytes() ([]byte, error)
	SetUUID(string)
	Position() protocol.AbsolutePosition
	Rotation() Rotation
}
