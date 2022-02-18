package generator

import (
	gameAPI "github.com/anchormc/anchor/src/api/game"
	"github.com/anchormc/anchor/src/impl/game"
)

func init() {
	GeneratorsMap["flat"] = NewFlatGenerator
}

type FlatGenerator struct {
	seed int64
}

func NewFlatGenerator(seed int64) gameAPI.Generator {
	return &FlatGenerator{
		seed: seed,
	}
}

func (g FlatGenerator) Seed() int64 {
	return g.seed
}

func (g FlatGenerator) GenerateChunk(x, z, height int64) (gameAPI.Chunk, error) {
	chunk := game.NewEmptyChunk(x, z, height)

	// Bedrock
	for bx := 0; bx < 16; bx++ {
		for bz := 0; bz < 16; bz++ {
			chunk.SetBlock(int64(bx), 0, int64(bz), 33)
		}
	}

	// Dirt
	for bx := 0; bx < 16; bx++ {
		for bz := 0; bz < 16; bz++ {
			for by := 1; by < 4; by++ {
				chunk.SetBlock(int64(bx), int64(by), int64(bz), 10)
			}
		}
	}

	// Grass
	for bx := 0; bx < 16; bx++ {
		for bz := 0; bz < 16; bz++ {
			chunk.SetBlock(int64(bx), 4, int64(bz), 8)
		}
	}

	return chunk, nil
}
