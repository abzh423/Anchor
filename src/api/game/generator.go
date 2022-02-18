package game

type Generator interface {
	Seed() int64
	GenerateChunk(x, y, z int64) (Chunk, error)
}
