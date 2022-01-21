package generators

import "github.com/golangminecraft/minecraft-server/src/api/world/chunk"

func init() {
	GeneratorsMap["flat"] = flatGenerator{}
}

type flatGenerator struct {
}

func (f flatGenerator) Initialize() error {
	return nil
}

func (f flatGenerator) GenerateChunk(x, z int64) (*chunk.Chunk, error) {
	c := chunk.NewChunk()

	var bx int64
	var by int64
	var bz int64

	for bx = 0; bx < 16; bx++ {
		for bz = 0; bz < 16; bz++ {
			for by = 0; by < 5; by++ {
				if by == 0 {
					c.SetBlock(bx, by, bz, chunk.PaletteItem{
						ID:         7,
						Name:       "minecraft:bedrock",
						Properties: make(map[string]interface{}),
					})
				} else if by >= 1 && by <= 3 {
					c.SetBlock(bx, by, bz, chunk.PaletteItem{
						ID:         3,
						Name:       "minecraft:dirt",
						Properties: make(map[string]interface{}),
					})
				} else if by == 4 {
					c.SetBlock(bx, by, bz, chunk.PaletteItem{
						ID:         2,
						Name:       "minecraft:grass_block",
						Properties: make(map[string]interface{}),
					})
				}
			}
		}
	}

	return c, nil
}

func (f flatGenerator) Close() error {
	return nil
}
