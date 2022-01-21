package generators

import "github.com/golangminecraft/minecraft-server/src/api/world/generator"

var GeneratorsMap map[string]generator.WorldGenerator = make(map[string]generator.WorldGenerator)
