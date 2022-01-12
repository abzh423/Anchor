package game

import "github.com/golangminecraft/minecraft-server/src/api/game/world"

type WorldManager struct {
	worlds map[string]world.World
}

func NewWorldManager() *WorldManager {
	return &WorldManager{
		worlds: make(map[string]world.World),
	}
}

func (w WorldManager) GetWorld(name string) world.World {
	return w.worlds[name]
}

func (w *WorldManager) NewWorld(world world.World) {
	w.worlds[world.Name()] = world
}

func (w WorldManager) Count() int {
	return len(w.worlds)
}

func (w WorldManager) Names() []string {
	names := make([]string, 0)

	for _, v := range w.worlds {
		names = append(names, v.Name())
	}

	return names
}

var _ world.WorldManager = &WorldManager{}
