package packets

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/golangminecraft/minecraft-server/src/enum"
	"github.com/golangminecraft/minecraft-server/src/protocol"
)

type HandshakePacket struct {
	ProtocolVersion int32
	ServerAddress   string
	ServerPort      uint16
	NextState       enum.ClientState
}

func (p HandshakePacket) Encode(w io.Writer) error {
	buf := &bytes.Buffer{}

	// Packet ID - varint
	if err := protocol.WriteVarInt(0x00, buf); err != nil {
		return err
	}

	// Protocol version - varint
	if err := protocol.WriteVarInt(p.ProtocolVersion, buf); err != nil {
		return err
	}

	// Server address - string
	if err := protocol.WriteString(p.ServerAddress, buf); err != nil {
		return err
	}

	// Server port - uint16
	if err := binary.Write(w, binary.BigEndian, p.ServerPort); err != nil {
		return err
	}

	// Next state - varint
	if err := protocol.WriteVarInt(int32(p.NextState), buf); err != nil {
		return err
	}

	if err := protocol.WriteVarInt(int32(buf.Len()), w); err != nil {
		return err
	}

	_, err := w.Write(buf.Bytes())

	return err
}

func (p *HandshakePacket) Decode(r io.Reader) error {
	protocolVersion, _, err := protocol.ReadVarInt(r)

	if err != nil {
		return err
	}

	serverAddress, err := protocol.ReadString(r)

	if err != nil {
		return err
	}

	var serverPort uint16

	if err := binary.Read(r, binary.BigEndian, &serverPort); err != nil {
		return err
	}

	nextState, _, err := protocol.ReadVarInt(r)

	if err != nil {
		return err
	}

	p.ProtocolVersion = protocolVersion
	p.ServerAddress = serverAddress
	p.ServerPort = serverPort
	p.NextState = enum.ClientState(nextState)

	return nil
}

var _ GenericPacket = &HandshakePacket{}
