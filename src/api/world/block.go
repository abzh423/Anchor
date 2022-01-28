package world

type BlockUpdate struct {
	X        int64
	Y        int64
	Z        int64
	OldValue Block
	NewValue Block
}

type Block struct {
	ID         int32                  `nbt:"ID"`
	Name       string                 `nbt:"Name"`
	Properties map[string]interface{} `nbt:"Properties"`
}
