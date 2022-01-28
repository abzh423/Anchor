package world

type NewWorldOptions struct {
	Folder       string
	Generator    WorldGenerator
	StoreOptions map[string]interface{}
}
