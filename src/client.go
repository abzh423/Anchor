package src

import (
	"io"
	"net"

	"github.com/golangminecraft/minecraft-server/src/enum"
)

type Client struct {
	Conn  net.Conn
	State enum.ClientState
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		Conn: conn,
	}
}

func (c Client) Process() error {
	for {
		packetID := make([]byte, 1)

		n, err := c.Conn.Read(packetID)

		if err != nil {
			return err
		}

		if n < 1 {
			return io.EOF
		}

		switch packetID[0] {
		case 0x00: // Handshake
			{
				switch c.State {
				case enum.ClientStateNone:
					{
						break
					}
				}
				break
			}
		}
	}
}

func (c Client) Close() error {
	return c.Conn.Close()
}
