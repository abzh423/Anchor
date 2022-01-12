package world

type Chunk interface {
	GetBlock(x, y, z int64) Block
	SetBlock(x, y, z int64, block Block)
}
