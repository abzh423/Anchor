package handlers

import (
	"io"

	"github.com/golangminecraft/minecraft-server/src/api/enum"
	proto "github.com/golangminecraft/minecraft-server/src/api/protocol"
	"github.com/golangminecraft/minecraft-server/src/api/server"
	"github.com/golangminecraft/minecraft-server/src/protocol"
)

type PingHandler struct{}

func (h PingHandler) PacketID() proto.VarInt {
	return 0x01
}

func (h PingHandler) Requirements(server server.Server, client server.Client) bool {
	return client.GetState() == enum.ClientStateStatus
}

func (h PingHandler) Execute(server server.Server, client server.Client, r io.Reader) error {
	var payload proto.Long

	if err := protocol.Unmarshal(r, &payload); err != nil {
		return err
	}

	packetData, err := protocol.Marshal(
		proto.VarInt(0x01),  // Pong (0x01)
		proto.Long(payload), // Unique ID
	)

	if err != nil {
		return err
	}

	return client.WritePacket(packetData)
}

var _ server.PacketHandler = &PingHandler{}
