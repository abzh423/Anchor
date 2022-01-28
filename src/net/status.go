package net

import (
	"encoding/json"
	"fmt"

	"github.com/anchormc/anchor/src/api"
	"github.com/anchormc/anchor/src/api/data"
	proto "github.com/anchormc/anchor/src/api/protocol"
)

func Status(server api.Server, client api.Client) error {
	requestPacket, err := client.ReadPacket()

	if err != nil {
		return err
	}

	if requestPacket.PacketID() != 0x00 {
		return fmt.Errorf("unexpected packet type (wanted: 0x00, received: 0x%02X)", requestPacket.PacketID())
	}

	favicon, err := server.Favicon()

	if err != nil {
		return err
	}

	data, err := json.Marshal(data.StatusResponse{
		Version: data.StatusResponseVersion{
			Name:     "1.18.1",
			Protocol: 757,
		},
		Players: data.StatusResponsePlayers{
			Online: server.OnlinePlayers(),
			Max:    server.MaxPlayers(),
			Sample: server.SamplePlayers(),
		},
		Description: server.MOTD(),
		Favicon:     favicon,
	})

	if err != nil {
		return err
	}

	statusPacket, err := proto.Marshal(
		proto.VarInt(0x00),
		proto.String(string(data)),
	)

	if err != nil {
		return err
	}

	if err = client.WritePacket(*statusPacket); err != nil {
		return err
	}

	pingPacket, err := client.ReadPacket()

	if err != nil {
		return err
	}

	if pingPacket.PacketID() != 0x01 {
		return fmt.Errorf("unexpected packet type (wanted: 0x01, received: 0x%02X)", requestPacket.PacketID())
	}

	var payload proto.Long

	if err = pingPacket.Unmarshal(&payload); err != nil {
		return err
	}

	pongPacket, err := proto.Marshal(
		proto.VarInt(0x01),
		payload,
	)

	if err != nil {
		return err
	}

	return client.WritePacket(*pongPacket)
}
