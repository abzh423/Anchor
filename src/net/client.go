package net

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"time"

	"github.com/golangminecraft/minecraft-server/src/api"
	"github.com/golangminecraft/minecraft-server/src/api/enum"
	proto "github.com/golangminecraft/minecraft-server/src/api/protocol"
	"github.com/golangminecraft/minecraft-server/src/components"
	"github.com/golangminecraft/minecraft-server/src/handlers"
	"github.com/google/uuid"
)

type Client struct {
	uuid                  uuid.UUID
	conn                  net.Conn
	reader                io.Reader
	writer                io.Writer
	player                api.Player
	compressionThreshold  int32
	lastSentKeepAliveID   int64
	lastSentKeepAliveTime time.Time
	latency               int64
}

func NewClient(conn net.Conn) (api.Client, error) {
	uuid, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	return &Client{
		uuid:                 uuid,
		conn:                 conn,
		reader:               conn,
		writer:               conn,
		player:               nil,
		compressionThreshold: -1,
	}, nil
}

func (c Client) ID() string {
	return c.uuid.String()
}

func (c Client) GetPlayer() api.Player {
	return c.player
}

func (c *Client) SetPlayer(player api.Player) {
	c.player = player
}

func (c Client) GetReader() io.Reader {
	return c.reader
}

func (c Client) GetWriter() io.Writer {
	return c.writer
}

func (c *Client) HandleConnection(server api.Server) error {
	handshake, err := Handshake(c)

	if err != nil {
		return err
	}

	defer (func() {
		for _, component := range components.Components {
			component.RemoveClient(c)
		}
	})()

	switch handshake.NextState {
	case enum.ClientStateStatus:
		{
			if err := Status(server, c); err != nil {
				log.Println(err)
			}
		}
	case enum.ClientStateLogin:
		{
			if err := Login(server, c); err != nil {
				log.Println(err)
			}

			for _, component := range components.Components {
				component.AddClient(c)
			}

			if err := Gameplay(server, c); err != nil {
				log.Println(err)
			}

			go c.StartKeepAlive(server.KeepAliveInterval())

			for server.Running() {
				packet, err := c.ReadPacket()

				if err != nil {
					return err
				}

				for _, handler := range handlers.Handlers {
					if handler.PacketID() != packet.PacketID() {
						continue
					}

					if err := handler.Execute(server, c, *packet); err != nil {
						log.Println(err)
					}

					break
				}
			}

			break
		}
	default:
		{
			log.Printf("received unknown next state value: %d\n", handshake.NextState)
		}
	}

	log.Printf("Client %s has disconnected\n", c.RemoteAddr())

	return c.Close()
}

func (c Client) Close() error {
	return c.conn.Close()
}

func (c Client) ReadPacket() (*proto.Packet, error) {
	var packetLength proto.VarInt

	if err := proto.Unmarshal(c.reader, &packetLength); err != nil {
		return nil, err
	}

	data := make([]byte, packetLength)

	if _, err := c.reader.Read(data); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(data)

	var packetID proto.VarInt

	if _, err := packetID.Decode(buf); err != nil {
		return nil, err
	}

	return proto.NewPacket(packetID, buf.Bytes()), nil
}

func (c Client) WritePacket(packet proto.Packet) error {
	packetBuffer := &bytes.Buffer{}

	if _, err := packet.PacketID().Encode(packetBuffer); err != nil {
		return err
	}

	if _, err := packetBuffer.Write(packet.Data()); err != nil {
		return err
	}

	if _, err := proto.VarInt(packetBuffer.Len()).Encode(c.writer); err != nil {
		return err
	}

	_, err := c.writer.Write(packetBuffer.Bytes())

	return err
}

func (c Client) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Client) SetCipher(encode, decode cipher.Stream) {
	c.reader = cipher.StreamReader{
		S: decode,
		R: c.conn,
	}

	c.writer = cipher.StreamWriter{
		S: encode,
		W: c.conn,
	}
}

func (c *Client) SetCompressionThreshold(threshold int32) {
	c.compressionThreshold = threshold
}

func (c *Client) StartKeepAlive(interval time.Duration) {
	for {
		payload := make([]byte, 8)

		if _, err := rand.Read(payload); err != nil {
			log.Println(err)

			break
		}

		c.lastSentKeepAliveID = int64(binary.BigEndian.Uint64(payload))

		keepAlivePacket, err := proto.Marshal(
			proto.VarInt(0x21),
			proto.Long(c.lastSentKeepAliveID),
		)

		if err != nil {
			log.Println(err)

			break
		}

		if err = c.WritePacket(*keepAlivePacket); err != nil {
			log.Println(err)

			break
		}

		c.lastSentKeepAliveTime = time.Now()

		<-time.NewTimer(interval).C
	}
}

func (c *Client) OnReceiveKeepAlive(payload int64) error {
	if payload != c.lastSentKeepAliveID {
		return errors.New("keep alive payload does not match last sent payload")
	}

	c.latency = int64(float64(time.Since(c.lastSentKeepAliveTime)) / float64(time.Millisecond))

	return nil
}

func (c *Client) SetLatency(latency int64) {
	c.latency = latency
}

func (c Client) Latency() int64 {
	return c.latency
}

var _ api.Client = &Client{}
