package components

import (
	"time"

	"github.com/anchormc/anchor/src/api"
	log "github.com/anchormc/anchor/src/api/logger"
	proto "github.com/anchormc/anchor/src/api/protocol"
)

func init() {
	Components = append(Components, &BroadcastLatency{
		clients: make(map[string]api.Client),
	})
}

type BroadcastLatency struct {
	isRunning bool
	server    api.Server
	clients   map[string]api.Client
}

func (c *BroadcastLatency) Initialize(server api.Server) error {
	c.server = server

	return nil
}

func (c *BroadcastLatency) Start() error {
	c.isRunning = true

	go c.Run()

	return nil
}

func (c *BroadcastLatency) Run() error {
	for c.isRunning {
		for _, client := range c.clients {
			packetArgs := []proto.DataTypeWriter{
				proto.VarInt(2),
				proto.VarInt(c.server.OnlinePlayers()),
			}

			for _, client := range c.server.Clients() {
				if client.GetPlayer() == nil {
					continue
				}

				packetArgs = append(packetArgs, client.GetPlayer().UUID())
				packetArgs = append(packetArgs, proto.VarInt(client.Latency()))
			}

			playerInfoPacket, err := proto.Marshal(
				proto.VarInt(0x36),
				packetArgs...,
			)

			if err != nil {
				return err
			}

			if err := client.WritePacket(*playerInfoPacket); err != nil {
				log.Error("broadcastlatencycomponent", err)
			}
		}

		<-time.NewTimer(time.Second * 3).C
	}

	return nil
}

func (c *BroadcastLatency) Stop() error {
	c.isRunning = false

	return nil
}

func (c *BroadcastLatency) AddClient(client api.Client) {
	c.clients[client.ID()] = client
}

func (c *BroadcastLatency) RemoveClient(client api.Client) {
	delete(c.clients, client.ID())
}
