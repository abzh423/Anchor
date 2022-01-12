package world

type Block interface {
	Identifier() string
	Properties() map[string]interface{}
}
