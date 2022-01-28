package world

import "io"

const (
	RegionSize      int64 = 32
	RegionBlockSize int64 = RegionSize * ChunkSize
)

type Region interface {
	X() int64
	Z() int64
	GetChunk(int64, int64) Chunk
	CreateChunk(int64, int64) Chunk
	Encode(io.Writer) error
	Decode(io.Reader) error
}
