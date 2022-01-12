package game

import "github.com/golangminecraft/minecraft-server/src/api/game/world"

type Chunk struct {
	x       int64
	z       int64
	blocks  []uint16
	palette []*world.Block
}

func NewChunk(x, z int64) *Chunk {
	return &Chunk{
		x:       x,
		z:       z,
		blocks:  make([]uint16, 16*16*256),
		palette: make([]*world.Block, 0),
	}
}

func (c Chunk) GetBlock(x, y, z int64) world.Block {
	paletteIndex := c.blocks[x+z*16+y*16*16]

	return *c.palette[paletteIndex]
}

func (c *Chunk) SetBlock(x, y, z int64, block world.Block) {
	blockIndex := x + z*16 + y*16*16

	c.blocks[blockIndex] = uint16(len(c.palette))

	c.palette = append(c.palette, &block)

paletteLoop:
	for k := range c.palette {
		for _, v := range c.blocks {
			if v == uint16(k) {
				continue paletteLoop
			}
		}

		oldPalette := c.palette
		c.palette = c.palette[0:k]
		c.palette = append(c.palette, oldPalette[k+1:]...)
	}
}

var _ world.Chunk = &Chunk{}
