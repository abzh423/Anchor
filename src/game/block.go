package game

import "github.com/golangminecraft/minecraft-server/src/api/game/world"

type Block struct {
	identifier string
	properties map[string]interface{}
}

func NewBlock(identifier string, properties map[string]interface{}) *Block {
	return &Block{
		identifier: identifier,
		properties: properties,
	}
}

func (p Block) Identifier() string {
	return p.identifier
}

func (p Block) Properties() map[string]interface{} {
	return p.properties
}

var _ world.Block = &Block{}
