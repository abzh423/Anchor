package networking

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
	"net"
	"reflect"
	"strings"

	"github.com/golangminecraft/minecraft-server/src/api/enum"
	"github.com/golangminecraft/minecraft-server/src/api/game"
	proto "github.com/golangminecraft/minecraft-server/src/api/protocol"
	"github.com/golangminecraft/minecraft-server/src/api/server"
	"github.com/golangminecraft/minecraft-server/src/protocol"
	"github.com/google/uuid"
)

type Client struct {
	id          string
	conn        net.Conn
	r           io.Reader
	w           io.Writer
	player      game.Player
	state       enum.ClientState
	verifyToken []byte
}

func NewClient(conn net.Conn) *Client {
	uuid, err := uuid.NewRandom()

	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		id:          strings.ToUpper(uuid.String()),
		conn:        conn,
		r:           conn,
		w:           conn,
		player:      nil,
		state:       enum.ClientStateNone,
		verifyToken: make([]byte, 0),
	}
}

func (c Client) ID() string {
	return c.id
}

func (c Client) GetReader() io.Reader {
	return c.r
}

func (c Client) GetWriter() io.Writer {
	return c.w
}

func (c *Client) SetCipher(encode, decode cipher.Stream) {
	c.r = cipher.StreamReader{
		S: decode,
		R: c.conn,
	}

	c.w = cipher.StreamWriter{
		S: encode,
		W: c.conn,
	}
}

func (c Client) GetPlayer() game.Player {
	return c.player
}

func (c *Client) SetPlayer(player game.Player) {
	c.player = player
}

func (c Client) GetState() enum.ClientState {
	return c.state
}

func (c *Client) SetState(state enum.ClientState) {
	c.state = state
}

func (c *Client) HandlePackets(server server.Server) error {
packetLoop:
	for {
		var packetLength proto.VarInt
		var packetType proto.VarInt

		if err := protocol.Unmarshal(c.r, &packetLength, &packetType); err != nil {
			return err
		}

		packetTypeLength := packetType.ByteLength()

		for _, handler := range server.GetPacketHandlers() {
			if handler.PacketID() != packetType || !handler.Requirements(server, c) {
				continue
			}

			log.Printf("Client %s: received packet (handler: %s, type: 0x%02X, length: %d)\n", c.id, reflect.TypeOf(handler).Name(), packetType, int64(packetLength)-packetTypeLength)

			if err := handler.Execute(server, c, c.r); err != nil {
				c.Disconnect()

				return err
			}

			continue packetLoop
		}

		data := make([]byte, int64(packetLength)-packetTypeLength)

		if _, err := c.conn.Read(data); err != nil {
			return err
		}

		log.Printf("Client %s: no packet handler for packet type 0x%02X (length: %d)\n", c.id, packetType, int64(packetLength)-packetTypeLength)
	}
}

func (c Client) GetConnection() net.Conn {
	return c.conn
}

func (c Client) WritePacket(data []byte) error {
	buffer := &bytes.Buffer{}

	if _, err := proto.VarInt(len(data)).Encode(buffer); err != nil {
		return err
	}

	if _, err := buffer.Write(data); err != nil {
		return err
	}

	_, err := c.w.Write(buffer.Bytes())

	log.Printf("Client %s: sent packet (length: %d)\n", c.id, len(data))

	return err
}

func (c Client) Disconnect() error {
	return c.conn.Close()
}

func (c *Client) GenerateVerifyToken() ([]byte, error) {
	c.verifyToken = make([]byte, 16)

	_, err := rand.Read(c.verifyToken)

	return c.verifyToken, err
}

func (c *Client) GetVerifyToken() []byte {
	return c.verifyToken
}

var _ server.Client = &Client{}
