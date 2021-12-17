package src

import (
	"bytes"
	"net"

	"github.com/golangminecraft/minecraft-server/src/enum"
	"github.com/golangminecraft/minecraft-server/src/protocol"
	"github.com/golangminecraft/minecraft-server/src/protocol/packets"
	"github.com/golangminecraft/minecraft-server/src/structure"
	"github.com/golangminecraft/minecraft-server/src/util"
)

type Client struct {
	UUID   string
	Server *Server
	Conn   net.Conn
	State  enum.ClientState
}

func NewClient(server *Server, conn net.Conn) *Client {
	uuid, _ := util.RandomUUID()

	return &Client{
		UUID:   uuid,
		Server: server,
		Conn:   conn,
	}
}

func (c Client) ReadPacket() (int32, *bytes.Buffer, error) {
	packetLength, _, err := protocol.ReadVarInt(c.Conn)

	if err != nil {
		return 0, nil, err
	}

	packetType, numRead, err := protocol.ReadVarInt(c.Conn)

	if err != nil {
		return packetType, nil, err
	}

	data := make([]byte, packetLength-int32(numRead))

	if _, err := c.Conn.Read(data); err != nil {
		return packetType, nil, err
	}

	return packetType, bytes.NewBuffer(data), nil
}

func (c *Client) Process() error {
	for {
		packetType, buffer, err := c.ReadPacket()

		if err != nil {
			return err
		}

		switch packetType {
		case 0x00:
			{
				switch c.State {
				case enum.ClientStateNone: // Handshake (0x00)
					{
						packet := packets.HandshakePacket{}

						if err := packet.Decode(buffer); err != nil {
							return err
						}

						c.State = packet.NextState

						break
					}
				case enum.ClientStateStatus: // Request (0x00)
					{
						onlinePlayers := 0

						for _, v := range c.Server.Clients {
							if v.State != enum.ClientStatePlay {
								continue
							}

							onlinePlayers++
						}

						packet := packets.ResponsePacket{
							JSONResponse: structure.JSONStatusResponse{
								Version: structure.StatusVersion{
									Name:     "1.18.1",
									Protocol: 757,
								},
								Players: structure.StatusPlayers{
									Online: onlinePlayers,
									Max:    int(c.Server.Config.MaxPlayers),
									Sample: make([]structure.StatusSamplePlayer, 0),
								},
								Description: c.Server.Config.MOTD,
								Favicon:     nil,
							},
						}

						if err := packet.Encode(c.Conn); err != nil {
							return err
						}

						break
					}
				}
				break
			}
		case 0x01:
			{
				pingPacket := packets.PingPacket{}

				if err := pingPacket.Decode(buffer); err != nil {
					return err
				}

				if err := pingPacket.Encode(c.Conn); err != nil {
					return err
				}

				break
			}
		}
	}
}

func (c Client) Close() error {
	return c.Conn.Close()
}
