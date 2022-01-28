package world

const (
	ChunkSize     int64 = 16
	ChunkHeight   int64 = 256
	SectionHeight int64 = 16
)

type Chunk interface {
	OnBlockUpdate() BlockUpdate
	BlockCount() int64
	GetSection(int) ChunkSection
	GetAllSections() []ChunkSection
	GetBlock(int64, int64, int64) *Block
	SetBlock(int64, int64, int64, Block)
}
