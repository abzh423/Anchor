package world

type WorldManager interface {
	LoadAllWorlds() error
	LoadWorld(string) error
	CreateWorld(string, WorldGenerator) (World, error)
	Worlds() []World
	GetWorld(string) (World, bool)
}
