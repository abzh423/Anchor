package enum

import "github.com/golangminecraft/minecraft-server/src/api/protocol"

var (
	PacketTypeHandshake protocol.VarInt = 0x00
)

type ClientState protocol.VarInt

var (
	ClientStateStatus ClientState = 1
	ClientStateLogin  ClientState = 2
)
