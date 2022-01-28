package world

type ChunkSection interface {
	OnBlockUpdate() BlockUpdate
	BlockCount() int16
	GetBlock(int64, int64, int64) *Block
	SetBlock(int64, int64, int64, Block)
	RawSectionData() ([]byte, error)
}
