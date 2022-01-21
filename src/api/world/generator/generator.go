package generator

import "github.com/golangminecraft/minecraft-server/src/api/world/chunk"

type WorldGenerator interface {
	Initialize() error
	GenerateChunk(int64, int64) (*chunk.Chunk, error)
	Close() error
}
