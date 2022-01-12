package protocol

type Packet interface {
	PacketID() VarInt
	Length() int
	Data() []byte
}
