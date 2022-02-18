package net

import (
	_ "embed"
	"math/rand"

	"github.com/Tnze/go-mc/nbt"
	"github.com/anchormc/anchor/src/api"
	"github.com/anchormc/protocol"
)

//go:embed DimensionCodec.snbt
var dimensionCodecSNBT []byte

//go:embed Dimension.snbt
var dimensionSNBT []byte

func Gameplay(server api.Server, client api.Client) error {
	if err := client.MarshalPacket(
		protocol.VarInt(0x26),
		protocol.Int(client.GetPlayer().EntityID()),
		protocol.Boolean(false),  // TODO hardcore mode
		protocol.UnsignedByte(1), // TODO gamemodes
		protocol.Byte(-1),
		protocol.VarInt(1), // TODO worlds
		protocol.Identifier("world"),
		protocol.NBT{Value: nbt.StringifiedMessage(dimensionCodecSNBT)},
		protocol.NBT{Value: nbt.StringifiedMessage(dimensionSNBT)},
		protocol.Identifier("minecraft:overworld"),
		protocol.Long(0), // TODO world seed
		protocol.VarInt(0),
		protocol.VarInt(12), // TODO server view distance
		protocol.VarInt(12), // TODO server simulation distance
		protocol.Boolean(false),
		protocol.Boolean(true),
		protocol.Boolean(false),
		protocol.Boolean(true), // TODO differentiate superflat worlds
	); err != nil {
		return err
	}

	if err := client.MarshalPacket(
		protocol.VarInt(0x0E),
		protocol.UnsignedByte(0), // TODO server difficulty
		protocol.Boolean(true),
	); err != nil {
		return err
	}

	var locale protocol.String
	var viewDistance protocol.Byte
	var chatMode protocol.VarInt
	var chatColors protocol.Boolean
	var displayedSkinParts protocol.UnsignedByte
	var mainHand protocol.VarInt
	var enableTextFiltering protocol.Boolean
	var allowServerListings protocol.Boolean

	if err := client.UnmarshalPacket(
		protocol.VarInt(0x05),
		&locale,
		&viewDistance,
		&chatMode,
		&chatColors,
		&displayedSkinParts,
		&mainHand,
		&enableTextFiltering,
		&allowServerListings,
	); err != nil {
		return err
	}

	position := client.GetPlayer().Position()

	if err := client.MarshalPacket(
		protocol.VarInt(0x38),
		protocol.Double(position.X),
		protocol.Double(position.Y),
		protocol.Double(position.Z),
		protocol.Float(0),
		protocol.Float(90),
		protocol.Byte(0),
		protocol.VarInt(rand.Int31()),
		protocol.Boolean(true),
	); err != nil {
		return err
	}

	return nil
}
