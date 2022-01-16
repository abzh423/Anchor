package world

type WorldGenerator interface {
	Initialize() error
	GenerateChunk(x, z int64) (*Chunk, error)
	Close() error
}
