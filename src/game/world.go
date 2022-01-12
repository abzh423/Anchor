package game

import (
	"fmt"

	"github.com/golangminecraft/minecraft-server/src/api/game/world"
)

type World struct {
	name      string
	chunks    map[string]world.Chunk
	generator world.WorldGenerator
	entities  []world.Entity
}

func NewWorld(name string, generator world.WorldGenerator) world.World {
	return &World{
		name:      name,
		chunks:    make(map[string]world.Chunk),
		generator: generator,
	}
}

func (w World) Name() string {
	return w.name
}

func (w World) NextEntityID() int {
	return len(w.entities)
}

func (w *World) AppendEntity(entity world.Entity) {
	w.entities = append(w.entities, entity)
}

func (w *World) GenerateChunk(x, z int64) {
	w.generator.GenerateChunk(x, z)
}

func (w World) GetBlock(x, y, z int64) world.Block {
	return w.GetChunk(x/16, z/16).GetBlock(x%16, y, z%16)
}

func (w World) SetBlock(x, y, z int64, block world.Block) {
	w.GetChunk(x/16, z/16).SetBlock(x%16, y, z%16, block)
}

func (w World) GetChunk(x, z int64) world.Chunk {
	chunkName := fmt.Sprintf("%d:%d", x/16, z/16)

	return w.chunks[chunkName]
}

func (w World) GetSeed() string {
	return w.generator.GetSeed()
}

var _ world.World = &World{}
