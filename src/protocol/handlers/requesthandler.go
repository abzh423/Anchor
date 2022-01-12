package handlers

import (
	"encoding/json"
	"io"

	"github.com/golangminecraft/minecraft-server/src/api/data"
	"github.com/golangminecraft/minecraft-server/src/api/enum"
	proto "github.com/golangminecraft/minecraft-server/src/api/protocol"
	"github.com/golangminecraft/minecraft-server/src/api/server"
	"github.com/golangminecraft/minecraft-server/src/protocol"
)

type RequestHandler struct{}

func (h RequestHandler) PacketID() proto.VarInt {
	return 0x00
}

func (h RequestHandler) Requirements(server server.Server, client server.Client) bool {
	return client.GetState() == enum.ClientStateStatus
}

func (h RequestHandler) Execute(server server.Server, client server.Client, r io.Reader) error {
	status := data.StatusResponse{
		Version: data.StatusResponseVersion{
			Name:     "1.18.1",
			Protocol: 757,
		},
		Players: data.StatusResponsePlayers{
			Online: server.OnlinePlayers(),
			Max:    server.GetConfig().MaxPlayers,
			Sample: make([]data.StatusResponseSamplePlayer, 0),
		},
		Description: data.Chat{
			Text: server.GetConfig().MOTD,
		},
		Favicon: nil,
	}

	data, err := json.Marshal(status)

	if err != nil {
		return err
	}

	packetData, err := protocol.Marshal(
		proto.VarInt(0x00), // Status Response (0x00)
		proto.String(data), // JSON Response
	)

	if err != nil {
		return err
	}

	return client.WritePacket(packetData)
}

var _ server.PacketHandler = &RequestHandler{}
