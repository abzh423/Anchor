package world

import "github.com/golangminecraft/minecraft-server/src/api/enum"

type Entity interface {
	Type() enum.EntityType
}
