package handlers

import (
	"github.com/golangminecraft/minecraft-server/src/api"
	proto "github.com/golangminecraft/minecraft-server/src/api/protocol"
)

func init() {
	Handlers = append(Handlers, &KeepAliveHandler{})
}

type KeepAliveHandler struct{}

func (k KeepAliveHandler) PacketID() proto.VarInt {
	return 0x0F
}

func (k KeepAliveHandler) Execute(server api.Server, client api.Client, packet proto.Packet) error {
	var keepAliveID proto.Long

	if err := packet.Unmarshal(&keepAliveID); err != nil {
		return err
	}

	return client.OnReceiveKeepAlive(int64(keepAliveID))
}
