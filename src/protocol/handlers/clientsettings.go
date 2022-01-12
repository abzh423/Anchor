package handlers

import (
	"io"

	"github.com/golangminecraft/minecraft-server/src/api/enum"
	proto "github.com/golangminecraft/minecraft-server/src/api/protocol"
	"github.com/golangminecraft/minecraft-server/src/api/server"
	"github.com/golangminecraft/minecraft-server/src/protocol"
)

type ClientSettingsHandler struct{}

func (h ClientSettingsHandler) PacketID() proto.VarInt {
	return 0x05
}

func (h ClientSettingsHandler) Requirements(server server.Server, client server.Client) bool {
	return client.GetState() == enum.ClientStatePlay
}

func (h ClientSettingsHandler) Execute(server server.Server, client server.Client, r io.Reader) error {
	var locale proto.String
	var viewDistance proto.Byte
	var chatMode proto.VarInt
	var chatColors proto.Boolean
	var displayedSkinParts proto.UnsignedByte
	var mainHand proto.VarInt
	var enableTextFiltering proto.Boolean
	var allowServerListings proto.Boolean

	if err := protocol.Unmarshal(r, &locale, &viewDistance, &chatMode, &chatColors, &displayedSkinParts, &mainHand, &enableTextFiltering, &allowServerListings); err != nil {
		return err
	}

	packetData, err := protocol.Marshal(
		proto.VarInt(0x48), // Held Item Change (0x48)
		proto.Byte(0),      // Slot ID
	)

	if err != nil {
		return err
	}

	if err := client.WritePacket(packetData); err != nil {
		return err
	}

	position := client.GetPlayer().Position()
	rotation := client.GetPlayer().Rotation()

	packetData, err = protocol.Marshal(
		proto.VarInt(0x38),          // Player Position and Look (0x38)
		proto.Double(position.X),    // X
		proto.Double(position.Y),    // Y
		proto.Double(position.Z),    // Z
		proto.Float(rotation.Yaw),   // Yaw
		proto.Float(rotation.Pitch), // Pitch
		proto.Byte(0),               // Flags (relative position & rotation)
		proto.VarInt(123456),        // Teleport ID
		proto.Boolean(true),         // Dismount Vehicle
	)

	if err != nil {
		return err
	}

	if err := client.WritePacket(packetData); err != nil {
		return err
	}

	// TODO send all players in this Player Info packet

	packetData, err = protocol.Marshal(
		proto.VarInt(0x36),                          // Player Info (0x36)
		proto.VarInt(0),                             // Action (add player)
		proto.VarInt(1),                             // Number of Players
		proto.UUID(client.GetPlayer().UUID()),       // UUID
		proto.String(client.GetPlayer().Username()), // Username
		proto.VarInt(0),                             // TODO properties
		proto.VarInt(0),                             // Gamemode
		proto.VarInt(10),                            // Ping (ms)
		proto.Boolean(false),                        // Has Display Name
	)

	if err != nil {
		return err
	}

	if err := client.WritePacket(packetData); err != nil {
		return err
	}

	packetData, err = protocol.Marshal(
		proto.VarInt(0x49), // Update View Position (0x49)
		proto.VarInt(client.GetPlayer().Position().X/16), // Chunk X
		proto.VarInt(client.GetPlayer().Position().Z/16), // Chunk Z
	)

	if err != nil {
		return err
	}

	if err := client.WritePacket(packetData); err != nil {
		return err
	}

	return nil
}

var _ server.PacketHandler = &ClientSettingsHandler{}
