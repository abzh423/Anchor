package world

import (
	"encoding/binary"
	"io"

	"github.com/Tnze/go-mc/nbt"
	"github.com/golangminecraft/minecraft-server/src/api/world/chunk"
)

const RegionSize int64 = 32

type Region struct {
	X      int64          `nbt:"x"`
	Z      int64          `nbt:"z"`
	chunks []*chunk.Chunk `nbt:"chunks"`
}

func LoadRegionFromData(data []byte) (*Region, error) {
	region := &Region{}

	if err := nbt.Unmarshal(data, region); err != nil {
		return nil, err
	}

	return region, nil
}

func NewRegion(x, z int64) *Region {
	chunks := make([]*chunk.Chunk, RegionSize*RegionSize)

	var cx int64
	var cz int64

	for cx = 0; cx < RegionSize; cx++ {
		for cz = 0; cz < RegionSize; cz++ {
			chunks[cx+cz*RegionSize] = chunk.NewChunk()
		}
	}

	return &Region{
		X:      x,
		Z:      z,
		chunks: chunks,
	}
}

func (r Region) GetChunk(x, z int64) *chunk.Chunk {
	return r.chunks[x+z*RegionSize]
}

func (r *Region) SetChunk(x, z int64, chunk *chunk.Chunk) {
	r.chunks[x+z*RegionSize] = chunk
}

func (r Region) IsDirty() bool {
	for _, chunk := range r.chunks {
		if chunk.IsDirty() {
			return true
		}
	}

	return false
}

func (r *Region) SetDirty(value bool) {
	for _, chunk := range r.chunks {
		chunk.SetDirty(value)
	}
}

func (s Region) Encode(w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, s.X); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, s.Z); err != nil {
		return err
	}

	for _, chunk := range s.chunks {
		if err := chunk.Encode(w); err != nil {
			return err
		}
	}

	return nil
}

func (s *Region) Decode(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, s.X); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, s.Z); err != nil {
		return err
	}

	for _, chunk := range s.chunks {
		if err := chunk.Decode(r); err != nil {
			return err
		}
	}

	return nil
}
