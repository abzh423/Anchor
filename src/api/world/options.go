package world

import "github.com/golangminecraft/minecraft-server/src/api/world/generator"

type NewWorldOptions struct {
	Folder       string
	Generator    generator.WorldGenerator
	StoreOptions map[string]interface{}
}
