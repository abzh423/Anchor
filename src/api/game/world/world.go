package world

type World interface {
	Name() string
	GetBlock(x, y, z int64) Block
	SetBlock(x, y, z int64, block Block)
	GetChunk(x, z int64) Chunk
	GenerateChunk(x, z int64)
	NextEntityID() int
	AppendEntity(Entity)
	GetSeed() string
}
