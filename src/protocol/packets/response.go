package packets

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/golangminecraft/minecraft-server/src/protocol"
	"github.com/golangminecraft/minecraft-server/src/structure"
)

type ResponsePacket struct {
	JSONResponse structure.JSONStatusResponse
}

func (p ResponsePacket) Encode(w io.Writer) error {
	buf := &bytes.Buffer{}

	// Packet ID - varint
	if err := protocol.WriteVarInt(0x00, buf); err != nil {
		return err
	}

	data, err := json.Marshal(p.JSONResponse)

	if err != nil {
		return err
	}

	// JSON Response
	if err := protocol.WriteString(string(data), buf); err != nil {
		return err
	}

	if err := protocol.WriteVarInt(int32(buf.Len()), w); err != nil {
		return err
	}

	_, err = w.Write(buf.Bytes())

	return err
}

func (p *ResponsePacket) Decode(r io.Reader) error {
	jsonResponse, err := protocol.ReadString(r)

	if err != nil {
		return err
	}

	response := structure.JSONStatusResponse{}

	if err := json.Unmarshal([]byte(jsonResponse), &response); err != nil {
		return err
	}

	p.JSONResponse = response

	return nil
}

var _ GenericPacket = &ResponsePacket{}
