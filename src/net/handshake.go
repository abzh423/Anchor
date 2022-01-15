package net

import (
	"fmt"

	"github.com/golangminecraft/minecraft-server/src/api"
	"github.com/golangminecraft/minecraft-server/src/api/enum"
	proto "github.com/golangminecraft/minecraft-server/src/api/protocol"
)

type HandshakeResponse struct {
	NextState enum.ClientState
}

func Handshake(client api.Client) (*HandshakeResponse, error) {
	packet, err := client.ReadPacket()

	if err != nil {
		return nil, err
	}

	if packet.PacketID() != enum.PacketTypeHandshake {
		return nil, fmt.Errorf("unexpected packet type (wanted: 0x%02X, received: 0x%02X)", enum.PacketTypeHandshake, packet.PacketID())
	}

	var protocolVersion proto.VarInt
	var serverHost proto.String
	var serverPort proto.UnsignedShort
	var nextState proto.VarInt

	if err = packet.Unmarshal(&protocolVersion, &serverHost, &serverPort, &nextState); err != nil {
		return nil, err
	}

	return &HandshakeResponse{
		NextState: enum.ClientState(nextState),
	}, nil
}
