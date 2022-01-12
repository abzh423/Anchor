package world

type WorldGenerator interface {
	GetSeed() string
	GenerateChunk(x, z int64) Chunk
}
