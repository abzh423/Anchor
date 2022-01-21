package chunk

import (
	"encoding/binary"
	"io"
)

const ChunkSize int = 16
const ChunkHeight int = 256
const SectionHeight int = 16

type Chunk struct {
	dirty    bool
	Sections []*ChunkSection `nbt:"sections"`
}

func NewChunk() *Chunk {
	sections := make([]*ChunkSection, 0)

	for i := 0; i < ChunkHeight; i += SectionHeight {
		sections = append(sections, NewChunkSection())
	}

	return &Chunk{
		dirty:    true,
		Sections: sections,
	}
}

func (c Chunk) BlockCount() (count int64) {
	for _, section := range c.Sections {
		count += int64(section.BlockCount())
	}

	return
}

func (c Chunk) IsDirty() bool {
	return c.dirty
}

func (c *Chunk) SetDirty(value bool) {
	c.dirty = value
}

func (c Chunk) GetSection(index int) *ChunkSection {
	return c.Sections[index]
}

func (c Chunk) GetBlock(x, y, z int64) *PaletteItem {
	return c.GetSection(int(y/int64(SectionHeight))).GetBlock(x%int64(ChunkSize), y%int64(ChunkSize), z%int64(ChunkSize))
}

func (c *Chunk) SetBlock(x, y, z int64, item PaletteItem) {
	c.GetSection(int(y/int64(SectionHeight))).SetBlock(x%int64(ChunkSize), y%int64(ChunkSize), z%int64(ChunkSize), item)
}

func (c Chunk) Encode(w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, uint32(len(c.Sections))); err != nil {
		return err
	}

	for _, section := range c.Sections {
		if err := section.Encode(w); err != nil {
			return err
		}
	}

	return nil
}

func (c *Chunk) Decode(r io.Reader) error {
	var sectionCount uint32

	if err := binary.Read(r, binary.BigEndian, &sectionCount); err != nil {
		return err
	}

	c.Sections = make([]*ChunkSection, sectionCount)

	for _, section := range c.Sections {
		if err := section.Decode(r); err != nil {
			return err
		}
	}

	return nil
}
