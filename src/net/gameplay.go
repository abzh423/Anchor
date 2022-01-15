package net

import (
	_ "embed"

	"github.com/Tnze/go-mc/nbt"
	"github.com/golangminecraft/minecraft-server/src/api"
	proto "github.com/golangminecraft/minecraft-server/src/api/protocol"
)

//go:embed DimensionCodec.snbt
var dimensionCodecSNBT string

//go:embed Dimension.snbt
var dimensionSNBT string

func Gameplay(server api.Server, client api.Client) error {
	joinGamePacket, err := proto.Marshal(
		proto.VarInt(0x26),
		proto.Int(client.GetPlayer().EntityID()),
		proto.Boolean(server.Hardcore()),
		proto.UnsignedByte(client.GetPlayer().Gamemode()),
		proto.Byte(-1),
		proto.VarInt(1), // TODO worlds and world count for Join Game Packet
		proto.Identifier("minecraft:overworld"),
		proto.NBT{Value: nbt.StringifiedMessage(dimensionCodecSNBT)},
		proto.NBT{Value: nbt.StringifiedMessage(dimensionSNBT)},
		proto.Identifier("minecraft:overworld"),
		proto.Long(123456), // TODO proper hashing of world seed
		proto.VarInt(server.MaxPlayers()),
		proto.VarInt(server.ViewDistance()),
		proto.VarInt(server.SimulationDistance()),
		proto.Boolean(false),
		proto.Boolean(true),
		proto.Boolean(false),
		proto.Boolean(false), // TODO set true if world is superflat
	)

	if err != nil {
		return err
	}

	if err := client.WritePacket(*joinGamePacket); err != nil {
		return err
	}

	difficultyPacket, err := proto.Marshal(
		proto.VarInt(0x0E),
		proto.UnsignedByte(server.Difficulty()),
		proto.Boolean(false),
	)

	if err != nil {
		return err
	}

	if err := client.WritePacket(*difficultyPacket); err != nil {
		return err
	}

	// TODO make sure plugin channel packets are ignored or properly disposed

	clientSettingsPacket, err := client.ReadPacket()

	if err != nil {
		return err
	}

	var locale proto.String
	var viewDistance proto.Byte
	var chatMode proto.VarInt
	var chatColors proto.Boolean
	var displayedSkinParts proto.UnsignedByte
	var mainHand proto.VarInt
	var enableTextFiltering proto.Boolean
	var allowServerListings proto.Boolean

	if err := clientSettingsPacket.Unmarshal(&locale, &viewDistance, &chatMode, &chatColors, &displayedSkinParts, &mainHand, &enableTextFiltering, &allowServerListings); err != nil {
		return err
	}

	position := client.GetPlayer().Position()
	rotation := client.GetPlayer().Rotation()

	playerPositionPacket, err := proto.Marshal(
		proto.VarInt(0x38),
		proto.Double(position.X),
		proto.Double(position.Y),
		proto.Double(position.Z),
		proto.Float(rotation.Yaw),
		proto.Float(rotation.Pitch),
		proto.Byte(0),
		proto.VarInt(0),
		proto.Boolean(true),
	)

	if err != nil {
		return err
	}

	if err := client.WritePacket(*playerPositionPacket); err != nil {
		return err
	}

	playerInfoArgs := []proto.DataTypeWriter{
		proto.VarInt(0), // Add player
		proto.VarInt(server.OnlinePlayers()),
	}

	for _, player := range server.Players() {
		playerInfoArgs = append(playerInfoArgs, player.UUID())
		playerInfoArgs = append(playerInfoArgs, proto.String(player.Username()))
		playerInfoArgs = append(playerInfoArgs, proto.VarInt(0))
		playerInfoArgs = append(playerInfoArgs, proto.VarInt(player.Gamemode()))
		playerInfoArgs = append(playerInfoArgs, proto.VarInt(50))
		playerInfoArgs = append(playerInfoArgs, proto.Boolean(false))
	}

	playerInfoPacket, err := proto.Marshal(
		proto.VarInt(0x36),
		playerInfoArgs...,
	)

	if err != nil {
		return err
	}

	return client.WritePacket(*playerInfoPacket)
}
