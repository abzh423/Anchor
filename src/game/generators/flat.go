package generators

import "github.com/golangminecraft/minecraft-server/src/api/world"

type flatGenerator struct {
}

func (f flatGenerator) Initialize() error {
	return nil
}

func (f flatGenerator) GenerateChunk(x, z int64) (*world.Chunk, error) {
	chunk := world.NewChunk(x, z)

	var bx int64
	var by int64
	var bz int64

	for bx = 0; bx < 16; bx++ {
		for bz = 0; bz < 16; bz++ {
			for by = 0; by < 5; by++ {
				if by == 0 {
					chunk.SetBlock(bx, by, bz, world.PaletteItem{
						Name:       "minecraft:bedrock",
						Properties: make(map[string]interface{}),
					})
				} else if by >= 1 && by <= 3 {
					chunk.SetBlock(bx, by, bz, world.PaletteItem{
						Name:       "minecraft:dirt",
						Properties: make(map[string]interface{}),
					})
				} else if by == 4 {
					chunk.SetBlock(bx, by, bz, world.PaletteItem{
						Name:       "minecraft:grass_block",
						Properties: make(map[string]interface{}),
					})
				}
			}
		}
	}

	return chunk, nil
}

func (f flatGenerator) Close() error {
	return nil
}

var FlatGenerator world.WorldGenerator = flatGenerator{}
