package world

type WorldStore interface {
	Initialize(string, map[string]interface{}) error
	GetChunk(int64, int64) (*Chunk, bool, error)
	SetChunk(int64, int64, Chunk) error
	Close() error
}
