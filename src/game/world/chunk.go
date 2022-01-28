package world

import "github.com/anchormc/anchor/src/api/world"

type Chunk struct {
	X               int64                  `nbt:"X"`
	Z               int64                  `nbt:"Y"`
	Sections        []world.ChunkSection   `nbt:"Sections"`
	blockUpdateChan chan world.BlockUpdate `nbt:"-"`
}

func NewChunk(x, z int64) world.Chunk {
	sections := make([]world.ChunkSection, 0)

	var i int64

	for i = 0; i < world.ChunkHeight; i += world.SectionHeight {
		sections = append(sections, NewChunkSection(i))
	}

	blockUpdateChan := make(chan world.BlockUpdate)

	for _, section := range sections {
		go (func() {
			for {
				blockUpdateChan <- section.OnBlockUpdate()
			}
		})()
	}

	return &Chunk{
		X:               x,
		Z:               z,
		Sections:        sections,
		blockUpdateChan: blockUpdateChan,
	}
}

func (c Chunk) OnBlockUpdate() world.BlockUpdate {
	return <-c.blockUpdateChan
}

func (c Chunk) BlockCount() (count int64) {
	for _, section := range c.Sections {
		count += int64(section.BlockCount())
	}

	return
}

func (c Chunk) GetSection(index int) world.ChunkSection {
	return c.Sections[index]
}

func (c Chunk) GetAllSections() []world.ChunkSection {
	return c.Sections
}

func (c Chunk) GetBlock(x, y, z int64) *world.Block {
	return c.GetSection(int(y/int64(world.SectionHeight))).GetBlock(x%int64(world.ChunkSize), y%int64(world.ChunkSize), z%int64(world.ChunkSize))
}

func (c *Chunk) SetBlock(x, y, z int64, item world.Block) {
	c.GetSection(int(y/int64(world.SectionHeight))).SetBlock(x%int64(world.ChunkSize), y%int64(world.ChunkSize), z%int64(world.ChunkSize), item)
}

var _ world.Chunk = &Chunk{}
