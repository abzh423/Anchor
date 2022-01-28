package api

import proto "github.com/anchormc/anchor/src/api/protocol"

type Handler interface {
	PacketID() proto.VarInt
	Execute(Server, Client, proto.Packet) error
}
