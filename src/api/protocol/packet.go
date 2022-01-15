package protocol

import (
	"bytes"
	"io"
)

type Packet struct {
	id   VarInt
	data []byte
}

func NewPacket(id VarInt, data []byte) *Packet {
	return &Packet{
		id:   id,
		data: data,
	}
}

func (p Packet) PacketID() VarInt {
	return p.id
}

func (p Packet) Length() int {
	return len(p.data)
}

func (p Packet) Data() []byte {
	return p.data
}

func Marshal(packetID VarInt, values ...DataTypeWriter) (*Packet, error) {
	buffer := &bytes.Buffer{}

	for _, value := range values {
		if _, err := value.Encode(buffer); err != nil {
			return nil, err
		}
	}

	return NewPacket(packetID, buffer.Bytes()), nil
}

func MarshalBytes(values ...DataTypeWriter) ([]byte, error) {
	buffer := &bytes.Buffer{}

	for _, value := range values {
		if _, err := value.Encode(buffer); err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

func (p Packet) Unmarshal(values ...DataTypeReader) error {
	r := bytes.NewReader(p.data)

	return Unmarshal(r, values...)
}

func Unmarshal(r io.Reader, values ...DataTypeReader) error {
	for _, value := range values {
		if _, err := value.Decode(r); err != nil {
			return err
		}
	}

	return nil
}
