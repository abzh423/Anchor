package handlers

import (
	"bytes"
	"errors"

	"github.com/golangminecraft/minecraft-server/src/api"
	"github.com/golangminecraft/minecraft-server/src/api/data"
	proto "github.com/golangminecraft/minecraft-server/src/api/protocol"
	"github.com/golangminecraft/minecraft-server/src/api/world/chunk"
)

func init() {
	Handlers = append(Handlers, &PositionHandler{})
}

type PositionHandler struct{}

func (k PositionHandler) PacketID() proto.VarInt {
	return 0x11
}

func (k PositionHandler) Execute(server api.Server, client api.Client, packet proto.Packet) error {
	player := client.GetPlayer()

	if player == nil {
		return errors.New("client sent a player movement packet, but is not logged in")
	}

	var x proto.Double
	var y proto.Double
	var z proto.Double
	var onGround proto.Boolean

	if err := packet.Unmarshal(&x, &y, &z, &onGround); err != nil {
		return err
	}

	player.SetPosition(proto.AbsolutePosition{
		X: float64(x),
		Y: float64(y),
		Z: float64(z),
	})

	return nil // TODO fix chunk sending

	chunkX := int64(x / proto.Double(chunk.ChunkSize))
	chunkZ := int64(x / proto.Double(chunk.ChunkSize))

	if player.HasSentChunk(chunkX, chunkZ) {
		return nil
	}

	chunk, err := server.DefaultWorld().GetChunkOrGenerate(chunkX, chunkZ)

	if err != nil {
		return err
	}

	sectionData := &bytes.Buffer{}

	for _, section := range chunk.Sections {
		data, err := section.RawSectionData()

		if err != nil {
			return err
		}

		if _, err := sectionData.Write(data); err != nil {
			return err
		}
	}

	chunkPacket, err := proto.Marshal(
		proto.VarInt(0x22), // Chunk Data and Update Light (0x22)
		proto.Int(chunkX),  // Chunk X
		proto.Int(chunkZ),  // Chunk Z
		proto.NBT{Value: data.Heightmap{MotionBlocking: []int64{0x0100804020100804}}}, // Heightmaps
		proto.ByteArray(sectionData.Bytes()),                                          // Section Data
		proto.VarInt(0),                                                               // Number of block entities
		proto.Boolean(true),                                                           // Trust Edges
	)

	if err != nil {
		return err
	}

	return client.WritePacket(*chunkPacket)
}
