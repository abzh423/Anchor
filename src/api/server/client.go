package server

import (
	"crypto/cipher"
	"io"

	"github.com/golangminecraft/minecraft-server/src/api/enum"
	"github.com/golangminecraft/minecraft-server/src/api/game"
)

type Client interface {
	ID() string
	GetReader() io.Reader
	GetWriter() io.Writer
	SetCipher(cipher.Stream, cipher.Stream)
	GetPlayer() game.Player
	SetPlayer(game.Player)
	GetState() enum.ClientState
	SetState(enum.ClientState)
	HandlePackets(Server) error
	WritePacket([]byte) error
	GenerateVerifyToken() ([]byte, error)
	GetVerifyToken() []byte
	Disconnect() error
}
