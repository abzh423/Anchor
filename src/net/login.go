package net

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/Tnze/go-mc/net/CFB8"
	"github.com/anchormc/anchor/src/api"
	log "github.com/anchormc/anchor/src/api/logger"
	proto "github.com/anchormc/anchor/src/api/protocol"
	"github.com/anchormc/anchor/src/game"
	"github.com/anchormc/anchor/src/util"
	"github.com/anchormc/anchor/src/util/rest"
)

func Login(server api.Server, client api.Client) error {
	loginPacket, err := client.ReadPacket()

	if err != nil {
		return err
	}

	var nickname proto.String

	if err := loginPacket.Unmarshal(&nickname); err != nil {
		return err
	}

	username := string(nickname)
	var uuid proto.UUID

	if server.OnlineMode() {
		verifyToken := make([]byte, 16)

		if _, err := rand.Read(verifyToken); err != nil {
			return err
		}

		privateKey, err := rsa.GenerateKey(rand.Reader, 1024)

		if err != nil {
			return err
		}

		publicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)

		if err != nil {
			return err
		}

		encryptionRequestPacket, err := proto.Marshal(
			proto.VarInt(0x01),
			proto.String(""),
			proto.ByteArray(publicKey),
			proto.ByteArray(verifyToken),
		)

		if err != nil {
			return err
		}

		if err = client.WritePacket(*encryptionRequestPacket); err != nil {
			return err
		}

		encryptionResponsePacket, err := client.ReadPacket()

		if err != nil {
			return err
		}

		var sharedSecret proto.ByteArray
		var encryptedVerifyToken proto.ByteArray

		if err = encryptionResponsePacket.Unmarshal(&sharedSecret, &encryptedVerifyToken); err != nil {
			return err
		}

		decodedSharedSecret, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, sharedSecret)

		if err != nil {
			return err
		}

		decodedVerifyToken, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedVerifyToken)

		if err != nil {
			return err
		}

		if !bytes.Equal(decodedVerifyToken, verifyToken) {
			return fmt.Errorf("decrypted verify token does not match server token")
		}

		block, err := aes.NewCipher(decodedSharedSecret)

		if err != nil {
			return err
		}

		client.SetCipher(CFB8.NewCFB8Encrypt(block, decodedSharedSecret), CFB8.NewCFB8Decrypt(block, decodedSharedSecret))

		hash := util.AuthDigest("", decodedSharedSecret, publicKey)

		response, err := rest.Authenticate(string(nickname), hash)

		if err != nil {
			return err
		}

		username = response.Name
		uuid = proto.UUID(response.ID)
	} else {
		uuid = proto.UUID(fmt.Sprintf("OfflinePlayer:%s", username))
	}

	loginSuccessPacket, err := proto.Marshal(
		proto.VarInt(0x02),     // Login Success (0x02)
		uuid,                   // UUID
		proto.String(username), // Username
	)

	if err != nil {
		return err
	}

	if err = client.WritePacket(*loginSuccessPacket); err != nil {
		return err
	}

	log.Infof("login", "%s (%s) has joined the game\n", username, uuid)

	client.SetPlayer(game.NewPlayer(server.NextEntityID(), username, uuid, server.DefaultGamemode()))

	return nil
}
