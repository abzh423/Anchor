package world

type WorldGenerator interface {
	Initialize() error
	GenerateChunk(Chunk) error
	Close() error
}
