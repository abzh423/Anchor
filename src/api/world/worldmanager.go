package world

import (
	"fmt"
	"os"
	"path"

	"github.com/golangminecraft/minecraft-server/src/game/generators"
)

type WorldManager struct {
	cwd    string
	worlds map[string]*World
}

func NewWorldManager(cwd string) *WorldManager {
	return &WorldManager{
		cwd:    path.Join(cwd, "worlds"),
		worlds: make(map[string]*World),
	}
}

func (w *WorldManager) LoadAllWorlds() error {
	folders, err := os.ReadDir(w.cwd)

	if err != nil {
		return err
	}

	for _, folder := range folders {
		if err = w.LoadWorld(path.Join(w.cwd, folder.Name())); err != nil {
			return err
		}
	}

	return nil
}

func (w *WorldManager) LoadWorld(folder string) error {
	world, err := LoadWorldFromFolder(folder)

	if err != nil {
		return err
	}

	w.worlds[world.Name()] = world

	return nil
}

func (w *WorldManager) CreateWorld(name string) (*World, error) {
	if _, ok := w.worlds[name]; ok {
		return nil, fmt.Errorf("attempted to create a new world, but one already exists with name: %s", name)
	}

	if err := os.Mkdir(path.Join(w.cwd, name), 0777); err != nil {
		return nil, err
	}

	world := NewWorld(name, WorldOptions{
		Folder:    path.Join(w.cwd, name),
		Generator: generators.GeneratorsMap["flat"],
		Height:    256,
	})

	w.worlds[name] = world

	if err := world.Save(); err != nil {
		return nil, err
	}

	return world, nil
}

func (w WorldManager) Worlds() []*World {
	worlds := make([]*World, 0)

	for _, world := range w.worlds {
		worlds = append(worlds, world)
	}

	return worlds
}

func (w WorldManager) GetWorld(name string) (*World, bool) {
	world, ok := w.worlds[name]

	return world, ok
}
