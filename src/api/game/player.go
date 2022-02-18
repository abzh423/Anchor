package game

import "github.com/anchormc/protocol"

type Player interface {
	EntityID() int64
	Username() string
	UUID() string
	Position() protocol.AbsolutePosition
}
