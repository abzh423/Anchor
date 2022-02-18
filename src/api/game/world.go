package game

import "github.com/Tnze/go-mc/data/block"

type World interface {
	Name() string
	Folder() string
	LoadAllRegions() error
	GetRegion(x, z int64) Region
	GetChunk(x, z int64) Chunk
	GetBlock(x, y, z int64) *block.Block
	GenerateChunk(x, z int64) error
}
