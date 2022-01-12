package generators

import (
	"github.com/golangminecraft/minecraft-server/src/api/game/world"
	"github.com/golangminecraft/minecraft-server/src/game"
)

type FlatWorldGenerator struct {
	seed string
}

func NewFlatWorldGenerator(seed string) *FlatWorldGenerator {
	return &FlatWorldGenerator{
		seed: seed,
	}
}

func (g *FlatWorldGenerator) GenerateChunk(x, z int64) world.Chunk {
	chunk := game.NewChunk(x, z)

	// Bedrock level
	for x := int64(0); x < 16; x++ {
		for z := int64(0); z < 16; z++ {
			chunk.SetBlock(x, 0, z, game.NewBlock("minecraft:bedrock", nil))
		}
	}

	// Stone level
	for x := int64(0); x < 16; x++ {
		for z := int64(0); z < 16; z++ {
			for y := int64(1); y < 4; y++ {
				chunk.SetBlock(x, y, z, game.NewBlock("minecraft:stone", nil))
			}
		}
	}

	// Grass level
	for x := int64(0); x < 16; x++ {
		for z := int64(0); z < 16; z++ {
			chunk.SetBlock(x, 4, z, game.NewBlock("minecraft:grass_block", nil))
		}
	}

	return chunk
}

func (g FlatWorldGenerator) GetSeed() string {
	return g.seed
}

var _ world.WorldGenerator = &FlatWorldGenerator{}
