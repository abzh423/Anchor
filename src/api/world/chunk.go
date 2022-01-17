package world

import (
	"encoding/binary"
	"io"
)

type Chunk struct {
	dirty   bool
	x       int64
	z       int64
	blocks  []uint16
	palette []PaletteItem
}

func NewChunk(x, z int64) *Chunk {
	return &Chunk{
		dirty:  true,
		x:      x,
		z:      z,
		blocks: make([]uint16, 16*16*256),
		palette: []PaletteItem{
			{
				Name:       "minecraft:air",
				Properties: make(map[string]interface{}),
			},
		},
	}
}

func (c Chunk) IsDirty() bool {
	return c.dirty
}

func (c *Chunk) SetDirty(value bool) {
	c.dirty = value
}

func (c Chunk) GetBlock(x, y, z int64) *PaletteItem {
	index := x + z*16 + y*16*16

	paletteIndex := c.blocks[index]

	return &c.palette[paletteIndex]
}

func (c *Chunk) SetBlock(x, y, z int64, item PaletteItem) {
	// TODO optimize palette items

	index := x + z*16 + y*16*16

	c.blocks[index] = uint16(len(c.palette))

	c.palette = append(c.palette, item)
	c.dirty = true
}

func (c Chunk) Encode(w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, c.x); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, c.z); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, uint16(len(c.blocks))); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, c.blocks); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, uint16(len(c.palette))); err != nil {
		return err
	}

	for _, item := range c.palette {
		if err := item.Encode(w); err != nil {
			return err
		}
	}

	return nil
}

func (c *Chunk) Decode(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &c.x); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &c.z); err != nil {
		return err
	}

	var blockLength uint16

	if err := binary.Read(r, binary.BigEndian, &blockLength); err != nil {
		return err
	}

	c.blocks = make([]uint16, blockLength)

	if err := binary.Read(r, binary.BigEndian, &c.blocks); err != nil {
		return err
	}

	var paletteLength uint16
	var i uint16

	for i = 0; i < paletteLength; i++ {
		item := PaletteItem{}

		if err := item.Decode(r); err != nil {
			return err
		}

		c.palette = append(c.palette, item)
	}

	return nil
}
