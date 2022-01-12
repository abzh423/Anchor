package protocol

import (
	"bytes"
	"io"

	"github.com/golangminecraft/minecraft-server/src/api/protocol"
)

type Packet struct {
	id   protocol.VarInt
	data []byte
}

func (p Packet) PacketID() protocol.VarInt {
	return p.id
}

func (p Packet) Length() int {
	return len(p.data)
}

func (p Packet) Data() []byte {
	return p.data
}

func Marshal(values ...protocol.DataTypeWriter) ([]byte, error) {
	buffer := &bytes.Buffer{}

	for _, value := range values {
		if _, err := value.Encode(buffer); err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

func Unmarshal(r io.Reader, values ...protocol.DataTypeReader) error {
	for _, value := range values {
		if _, err := value.Decode(r); err != nil {
			return err
		}
	}

	return nil
}
