package handlers

import (
	"io"

	"github.com/golangminecraft/minecraft-server/src/api/enum"
	proto "github.com/golangminecraft/minecraft-server/src/api/protocol"
	"github.com/golangminecraft/minecraft-server/src/api/server"
	"github.com/golangminecraft/minecraft-server/src/protocol"
)

type HandshakeHandler struct{}

func (h HandshakeHandler) PacketID() proto.VarInt {
	return 0x00
}

func (h HandshakeHandler) Requirements(server server.Server, client server.Client) bool {
	return client.GetState() == enum.ClientStateNone
}

func (h HandshakeHandler) Execute(server server.Server, client server.Client, r io.Reader) error {
	var protocolVersion proto.VarInt
	var serverHost proto.String
	var serverPort proto.UnsignedShort
	var nextState proto.VarInt

	if err := protocol.Unmarshal(r, &protocolVersion, &serverHost, &serverPort, &nextState); err != nil {
		return err
	}

	client.SetState(enum.ClientState(nextState))

	return nil
}

var _ server.PacketHandler = &HandshakeHandler{}
