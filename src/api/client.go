package api

import (
	"crypto/cipher"
	"io"
	"net"
	"time"

	"github.com/golangminecraft/minecraft-server/src/api/protocol"
)

type Client interface {
	ID() string
	GetPlayer() Player
	SetPlayer(Player)
	GetReader() io.Reader
	GetWriter() io.Writer
	ReadPacket() (*protocol.Packet, error)
	WritePacket(protocol.Packet) error
	HandleConnection(Server) error
	RemoteAddr() net.Addr
	SetCipher(encode, decode cipher.Stream)
	SetCompressionThreshold(int32)
	StartKeepAlive(time.Duration)
	OnReceiveKeepAlive(int64) error
	Latency() int64
	Close() error
}
