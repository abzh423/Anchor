package server

import (
	"io"

	"github.com/golangminecraft/minecraft-server/src/api/protocol"
)

type PacketHandler interface {
	PacketID() protocol.VarInt
	Requirements(Server, Client) bool
	Execute(Server, Client, io.Reader) error
}
