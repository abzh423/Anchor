package world

type WorldManager interface {
	Count() int
	Names() []string
	GetWorld(name string) World
	NewWorld(World)
}
