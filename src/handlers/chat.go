package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/anchormc/anchor/src/api"
	proto "github.com/anchormc/anchor/src/api/protocol"
)

func init() {
	Handlers = append(Handlers, &ChatHandler{})
}

type ChatHandler struct{}

func (k ChatHandler) PacketID() proto.VarInt {
	return 0x03
}

func (k ChatHandler) Execute(server api.Server, client api.Client, packet proto.Packet) error {
	var message proto.String

	if err := packet.Unmarshal(&message); err != nil {
		return err
	}

	if strings.HasPrefix(string(message), "/") {
		// TODO client commands

		packetData, err := proto.Marshal(
			proto.VarInt(0x0F),
			proto.Chat{
				Text:  fmt.Sprintf("Unknown command: %s", message),
				Color: "red",
			},
			proto.Byte(1),
			client.GetPlayer().UUID(),
		)

		if err != nil {
			return err
		}

		if err := client.WritePacket(*packetData); err != nil {
			return err
		}

		return errors.New("client attempted to run command, not implemented")
	}

	player := client.GetPlayer()

	if player == nil {
		return errors.New("client sent a chat message, but is not logged in")
	}

	packetData, err := proto.Marshal(
		proto.VarInt(0x0F),
		proto.Chat{
			Text: fmt.Sprintf("%s: %s", player.Username(), message),
		},
		proto.Byte(0),
		client.GetPlayer().UUID(),
	)

	if err != nil {
		return err
	}

	for _, c := range server.Clients() {
		player := c.GetPlayer()

		if player == nil {
			continue
		}

		if err := c.WritePacket(*packetData); err != nil {
			return err
		}
	}

	return nil
}
