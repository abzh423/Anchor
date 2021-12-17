package packets

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/golangminecraft/minecraft-server/src/protocol"
)

type PingPacket struct {
	Time int64
}

func (p PingPacket) Encode(w io.Writer) error {
	buf := &bytes.Buffer{}

	// Packet ID - varint
	if err := protocol.WriteVarInt(0x01, buf); err != nil {
		return err
	}

	// Time - int64
	if err := binary.Write(buf, binary.BigEndian, p.Time); err != nil {
		return err
	}

	if err := protocol.WriteVarInt(int32(buf.Len()), w); err != nil {
		return err
	}

	_, err := w.Write(buf.Bytes())

	return err
}

func (p *PingPacket) Decode(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &p.Time)
}

var _ GenericPacket = &ResponsePacket{}
