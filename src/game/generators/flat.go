package generators

import (
	"github.com/anchormc/anchor/src/api/world"
)

func init() {
	GeneratorsMap["flat"] = FlatGenerator{}
}

type FlatGenerator struct {
}

func (f FlatGenerator) Initialize() error {
	return nil
}

func (f FlatGenerator) GenerateChunk(chunk world.Chunk) error {
	var bx int64
	var by int64
	var bz int64

	for bx = 0; bx < 16; bx++ {
		for bz = 0; bz < 16; bz++ {
			for by = 0; by < 5; by++ {
				if by == 0 {
					chunk.SetBlock(bx, by, bz, world.Block{
						ID:         7,
						Name:       "minecraft:bedrock",
						Properties: make(map[string]interface{}),
					})
				} else if by >= 1 && by <= 3 {
					chunk.SetBlock(bx, by, bz, world.Block{
						ID:         3,
						Name:       "minecraft:dirt",
						Properties: make(map[string]interface{}),
					})
				} else if by == 4 {
					chunk.SetBlock(bx, by, bz, world.Block{
						ID:         2,
						Name:       "minecraft:grass_block",
						Properties: make(map[string]interface{}),
					})
				}
			}
		}
	}

	return nil
}

func (f FlatGenerator) Close() error {
	return nil
}
