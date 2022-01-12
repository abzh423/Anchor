package handlers

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	_ "embed"
	"fmt"
	"io"

	"github.com/Tnze/go-mc/nbt"
	"github.com/Tnze/go-mc/net/CFB8"
	"github.com/golangminecraft/minecraft-server/src/api/enum"
	proto "github.com/golangminecraft/minecraft-server/src/api/protocol"
	"github.com/golangminecraft/minecraft-server/src/api/server"
	"github.com/golangminecraft/minecraft-server/src/protocol"
	"github.com/golangminecraft/minecraft-server/src/rest"
	"github.com/golangminecraft/minecraft-server/src/util"
)

type EncryptionResponseHandler struct{}

func (h EncryptionResponseHandler) PacketID() proto.VarInt {
	return 0x01
}

func (h EncryptionResponseHandler) Requirements(server server.Server, client server.Client) bool {
	return client.GetState() == enum.ClientStatePlay
}

//go:embed DimensionCodec.snbt
var dimensionCodecSNBT string

//go:embed Dimension.snbt
var dimensionSNBT string

func (h EncryptionResponseHandler) Execute(server server.Server, client server.Client, r io.Reader) error {
	var sharedSecret proto.ByteArray
	var verifyToken proto.ByteArray

	if err := protocol.Unmarshal(r, &sharedSecret, &verifyToken); err != nil {
		return err
	}

	privateKey := server.GetPrivateKey()

	publicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)

	if err != nil {
		return err
	}

	decodedSharedSecret, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, sharedSecret)

	if err != nil {
		return err
	}

	decodedVerifyToken, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, verifyToken)

	if err != nil {
		return err
	}

	if !bytes.Equal(decodedVerifyToken, client.GetVerifyToken()) {
		return fmt.Errorf("decrypted verify token does not match server token")
	}

	block, err := aes.NewCipher(decodedSharedSecret)

	if err != nil {
		return err
	}

	client.SetCipher(CFB8.NewCFB8Encrypt(block, decodedSharedSecret), CFB8.NewCFB8Decrypt(block, decodedSharedSecret))

	hash := util.AuthDigest("", decodedSharedSecret, publicKey)

	response, err := rest.Authenticate(client.GetPlayer().Username(), hash)

	if err != nil {
		return err
	}

	client.GetPlayer().SetUUID(response.ID)

	packetData, err := protocol.Marshal(
		proto.VarInt(0x02),          // Login Success (0x02)
		proto.UUID(response.ID),     // UUID
		proto.String(response.Name), // Username
	)

	if err != nil {
		return err
	}

	if err = client.WritePacket(packetData); err != nil {
		return err
	}

	packetFields := []proto.DataTypeWriter{
		proto.VarInt(0x26),                             // Join Game (0x02)
		proto.Int(client.GetPlayer().EntityID()),       // Entity ID
		proto.Boolean(false),                           // Hardcore
		proto.UnsignedByte(0),                          // Gamemode
		proto.Byte(-1),                                 // Previous Gamemode
		proto.VarInt(server.GetWorldManager().Count()), // World Count
	}

	for _, v := range server.GetWorldManager().Names() {
		packetFields = append(packetFields, proto.Identifier(v)) // Dimension Names
	}

	packetFields = append(packetFields, proto.NBT{Value: nbt.StringifiedMessage(dimensionCodecSNBT)})
	packetFields = append(packetFields, proto.NBT{Value: nbt.StringifiedMessage(dimensionSNBT)})
	packetFields = append(packetFields, proto.String("minecraft:overworld")) // World Name
	packetFields = append(packetFields, proto.Long(123456))                  // Hashed Seed
	packetFields = append(packetFields, proto.VarInt(0))                     // Max Players
	packetFields = append(packetFields, proto.VarInt(16))                    // View Distance
	packetFields = append(packetFields, proto.VarInt(10))                    // Simulation Distance
	packetFields = append(packetFields, proto.Boolean(false))                // Reduce Debug
	packetFields = append(packetFields, proto.Boolean(true))                 // Respawn Screen
	packetFields = append(packetFields, proto.Boolean(false))                // Debug Mode
	packetFields = append(packetFields, proto.Boolean(false))                // Flat World

	packetData, err = protocol.Marshal(packetFields...)

	if err != nil {
		return err
	}

	if err := client.WritePacket(packetData); err != nil {
		return err
	}

	packetData, err = protocol.Marshal(
		proto.VarInt(0x0E),    // Server Difficulty (0x0E)
		proto.UnsignedByte(0), // Difficulty ID
		proto.Boolean(true),   // Difficulty Locked
	)

	if err != nil {
		return err
	}

	return client.WritePacket(packetData)
}

var _ server.PacketHandler = &EncryptionResponseHandler{}
