package net

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	random "math/rand"

	"github.com/Tnze/go-mc/net/CFB8"
	"github.com/anchormc/anchor/src/api"
	"github.com/anchormc/anchor/src/api/util"
	"github.com/anchormc/anchor/src/impl/game"
	"github.com/anchormc/protocol"
	UUID "github.com/google/uuid"
)

func Login(server api.Server, client api.Client) error {
	var nickname protocol.String

	if err := client.UnmarshalPacket(
		protocol.VarInt(0x00),
		&nickname,
	); err != nil {
		return err
	}

	var username string
	var uuid protocol.UUID

	if server.GetConfig().OnlineMode {
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

		if err = client.MarshalPacket(
			protocol.VarInt(0x01),
			protocol.String(""),
			protocol.ByteArray(publicKey),
			protocol.ByteArray(verifyToken),
		); err != nil {
			return err
		}

		var sharedSecret protocol.ByteArray
		var encryptedVerifyToken protocol.ByteArray

		if err = client.UnmarshalPacket(
			protocol.VarInt(0x01),
			&sharedSecret,
			&encryptedVerifyToken,
		); err != nil {
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

		response, err := util.Authenticate(string(nickname), hash)

		if err != nil {
			return err
		}

		username = response.Name
		uuid = protocol.UUID(response.ID)
	} else {
		username = "OfflinePlayer"
		uuid = protocol.UUID(UUID.NewString())
	}

	if err := client.MarshalPacket(
		protocol.VarInt(0x02),
		uuid,
		protocol.String(username),
	); err != nil {
		return err
	}

	client.SetPlayer(game.NewPlayer(random.Int63(), username, string(uuid), protocol.AbsolutePosition{}))

	return nil
}
