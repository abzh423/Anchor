package generator

import "github.com/anchormc/anchor/src/api/game"

var GeneratorsMap map[string]func(int64) game.Generator = make(map[string]func(int64) game.Generator)
