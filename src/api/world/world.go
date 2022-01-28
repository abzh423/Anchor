package world

type World interface {
	Name() string
	Initialize() error
	GetRegion(int64, int64) (Region, bool, error)
	GetChunk(int64, int64) (Chunk, bool, error)
	GetChunkOrGenerate(int64, int64) (Chunk, error)
	GetBlock(int64, int64, int64) (*Block, bool, error)
	GenerateChunk(int64, int64) error
	CreateRegion(int64, int64) Region
	Save() error
	Close() error
}
