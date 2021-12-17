package packets

import "io"

type GenericPacket interface {
	Encode(w io.Writer) error
	Decode(r io.Reader) error
}
