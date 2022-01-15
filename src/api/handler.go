package api

import proto "github.com/golangminecraft/minecraft-server/src/api/protocol"

type Handler interface {
	PacketID() proto.VarInt
	Execute(Server, Client, proto.Packet) error
}
