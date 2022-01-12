package handlers

import (
	"crypto/x509"
	"io"

	"github.com/golangminecraft/minecraft-server/src/api/enum"
	proto "github.com/golangminecraft/minecraft-server/src/api/protocol"
	"github.com/golangminecraft/minecraft-server/src/api/server"
	"github.com/golangminecraft/minecraft-server/src/game"
	"github.com/golangminecraft/minecraft-server/src/protocol"
)

type LoginStartHandler struct{}

func (h LoginStartHandler) PacketID() proto.VarInt {
	return 0x00
}

func (h LoginStartHandler) Requirements(server server.Server, client server.Client) bool {
	return client.GetState() == enum.ClientStatePlay
}

func (h LoginStartHandler) Execute(server server.Server, client server.Client, r io.Reader) error {
	var username proto.String

	if err := protocol.Unmarshal(r, &username); err != nil {
		return err
	}

	world := server.GetWorldManager().GetWorld("world")

	player := game.NewPlayer(world.NextEntityID(), string(username))
	client.SetPlayer(player)

	world.AppendEntity(player)

	verifyToken, err := client.GenerateVerifyToken()

	if err != nil {
		return err
	}

	privateKey := server.GetPrivateKey()

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)

	if err != nil {
		return err
	}

	packetData, err := protocol.Marshal(
		proto.VarInt(0x01),              // Encryption Request (0x01)
		proto.String(""),                // Server ID
		proto.ByteArray(publicKeyBytes), // Public Key
		proto.ByteArray(verifyToken),    // Verify Token
	)

	if err != nil {
		return err
	}

	return client.WritePacket(packetData)
}

var _ server.PacketHandler = &LoginStartHandler{}
