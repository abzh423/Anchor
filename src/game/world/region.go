package world

import (
	"io"

	"github.com/Tnze/go-mc/nbt"
	"github.com/anchormc/anchor/src/api/world"
)

type Region struct {
	XCoord int64         `nbt:"X"`
	ZCoord int64         `nbt:"Z"`
	Chunks []world.Chunk `nbt:"Chunks"`
}

func NewRegion(x, z int64) *Region {
	chunks := make([]world.Chunk, world.RegionSize*world.RegionSize)

	var cx int64
	var cz int64

	for cx = 0; cx < world.RegionSize; cx++ {
		for cz = 0; cz < world.RegionSize; cz++ {
			chunks[cx+cz*world.RegionSize] = NewChunk(cx+x*world.RegionSize, cz+z*world.RegionSize)
		}
	}

	return &Region{
		XCoord: x,
		ZCoord: z,
		Chunks: chunks,
	}
}

func (r Region) X() int64 {
	return r.XCoord
}

func (r Region) Z() int64 {
	return r.ZCoord
}

func (r Region) GetChunk(x, z int64) world.Chunk {
	return r.Chunks[x+z*world.RegionSize]
}

func (r *Region) CreateChunk(x, z int64) world.Chunk {
	chunk := NewChunk(x, z)

	r.Chunks[x+z*world.RegionSize] = chunk

	return chunk
}

func (r Region) Encode(writer io.Writer) error {
	data, err := nbt.Marshal(r)

	if err != nil {
		return err
	}

	_, err = writer.Write(data)

	return err
}

func (r *Region) Decode(reader io.Reader) error {
	decoder := nbt.NewDecoder(reader)

	_, err := decoder.Decode(r)

	return err
}

var _ world.Region = &Region{}
