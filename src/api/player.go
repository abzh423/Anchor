package api

import (
	"github.com/anchormc/anchor/src/api/enum"
	proto "github.com/anchormc/anchor/src/api/protocol"
)

type Player interface {
	EntityID() int64
	Username() string
	UUID() proto.UUID
	Gamemode() enum.Gamemode
	Position() proto.AbsolutePosition
	SetPosition(proto.AbsolutePosition)
	Chunk() proto.RelativePosition
	SetSentChunk(int64, int64)
	HasSentChunk(int64, int64) bool
	Rotation() Rotation
	SetRotation(Rotation)
}
