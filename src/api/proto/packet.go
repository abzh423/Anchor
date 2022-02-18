package proto

import (
	"bytes"

	"github.com/anchormc/protocol"
)

type Packet struct {
	ID     protocol.VarInt
	Buffer *bytes.Buffer
}

func NewPacket(id protocol.VarInt, data []byte) *Packet {
	return &Packet{
		ID:     id,
		Buffer: bytes.NewBuffer(data),
	}
}

func (p Packet) Read(values ...protocol.DataTypeReader) error {
	return protocol.Unmarshal(p.Buffer, values...)
}

func (p Packet) Write(values ...protocol.DataTypeWriter) error {
	data, err := protocol.Marshal(values...)

	if err != nil {
		return err
	}

	_, err = p.Buffer.Write(data)

	return err
}

func (p Packet) Length() int {
	return p.Buffer.Len()
}
