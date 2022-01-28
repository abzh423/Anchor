package handlers

import (
	"errors"

	"github.com/anchormc/anchor/src/api"
	proto "github.com/anchormc/anchor/src/api/protocol"
)

func init() {
	Handlers = append(Handlers, &RotationHandler{})
}

type RotationHandler struct{}

func (k RotationHandler) PacketID() proto.VarInt {
	return 0x13
}

func (k RotationHandler) Execute(server api.Server, client api.Client, packet proto.Packet) error {
	player := client.GetPlayer()

	if player == nil {
		return errors.New("client sent a player rotation packet, but is not logged in")
	}

	var yaw proto.Float
	var pitch proto.Float
	var onGround proto.Boolean

	if err := packet.Unmarshal(&yaw, &pitch, &onGround); err != nil {
		return err
	}

	player.SetRotation(api.Rotation{
		Yaw:   float64(yaw),
		Pitch: float64(pitch),
	})

	return nil
}
