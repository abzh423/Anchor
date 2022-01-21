package world

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/Tnze/go-mc/nbt"
	"github.com/golangminecraft/minecraft-server/src/api/world/chunk"
	"github.com/golangminecraft/minecraft-server/src/api/world/generator"
	"github.com/golangminecraft/minecraft-server/src/game/generators"
)

type World struct {
	cwd       string
	meta      WorldMeta
	generator generator.WorldGenerator
	regions   map[string]*Region
}

type WorldOptions struct {
	Folder    string
	Generator generator.WorldGenerator
	Height    int
}

type WorldMeta struct {
	Name      string `nbt:"name"`
	Generator string `nbt:"generator"`
	Height    int32  `nbt:"height"`
}

func LoadWorldFromFolder(cwd string) (*World, error) {
	_, err := os.Stat(cwd)

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path.Join(cwd, "level.dat"))

	if err != nil {
		return nil, err
	}

	meta := WorldMeta{}

	if err := nbt.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	generator, ok := generators.GeneratorsMap[meta.Generator]

	if !ok {
		return nil, fmt.Errorf("unknown generator for world '%s': %s", meta.Name, meta.Generator)
	}

	return &World{
		cwd:       cwd,
		meta:      meta,
		generator: generator,
		regions:   make(map[string]*Region),
	}, nil
}

func NewWorld(name string, options WorldOptions) *World {
	return &World{
		cwd: options.Folder,
		meta: WorldMeta{
			Name:      name,
			Generator: "flat",
			Height:    256,
		},
		generator: options.Generator,
		regions:   make(map[string]*Region),
	}
}

func (w World) Name() string {
	return w.meta.Name
}

func (w World) Initialize() error {
	if err := os.Mkdir(path.Join(w.cwd, "region"), 0777); err != nil {
		return err
	}

	return nil
}

func (w *World) GetRegion(x, z int64) (*Region, bool, error) {
	region, ok := w.regions[fmt.Sprintf("%d:%d", x, z)]

	if ok {
		return region, true, nil
	}

	_, err := os.Stat(path.Join(w.cwd, "region", fmt.Sprintf("region.%d.%d.bin", x, z)))

	if err != nil {
		return nil, false, nil
	}

	data, err := ioutil.ReadFile(path.Join(w.cwd, "region", fmt.Sprintf("region.%d.%d.bin", x, z)))

	if err != nil {
		return nil, false, err
	}

	region, err = LoadRegionFromData(data)

	if err != nil {
		return nil, true, err
	}

	w.regions[fmt.Sprintf("%d:%d", x, z)] = region

	return region, true, nil
}

func (w *World) GetChunk(x, z int64) (*chunk.Chunk, bool, error) {
	region, ok, err := w.GetRegion(x/RegionSize, z/RegionSize)

	if err != nil || !ok {
		return nil, ok, err
	}

	return region.GetChunk(x%RegionSize, z%RegionSize), true, nil
}

func (w *World) GetChunkOrGenerate(x, z int64) (*chunk.Chunk, error) {
	region, ok, err := w.GetRegion(x/RegionSize, z/RegionSize)

	if err != nil {
		return nil, err
	}

	if !ok {
		region = w.CreateRegion(x, z)

		if err := w.GenerateChunk(x, z); err != nil {
			return nil, err
		}
	}

	return region.GetChunk(x%RegionSize, z%RegionSize), nil
}

func (w *World) GetBlock(x, y, z int64) (*chunk.PaletteItem, bool, error) {
	c, ok, err := w.GetChunk(x/int64(chunk.ChunkSize), z/int64(chunk.ChunkSize))

	if err != nil || !ok {
		return nil, ok, err
	}

	return c.GetBlock(x%int64(chunk.ChunkSize), y, z%int64(chunk.ChunkSize)), true, nil
}

func (w *World) GenerateChunk(x, z int64) error {
	log.Printf("Generating chunk (%d,%d)\n", x, z)

	region, ok, err := w.GetRegion(x/RegionSize, z/RegionSize)

	if err != nil {
		return err
	}

	if !ok {
		region = w.CreateRegion(x/RegionSize, z/RegionSize)
	}

	chunk, err := w.generator.GenerateChunk(x, z)

	if err != nil {
		return err
	}

	region.SetChunk(x%RegionSize, z%RegionSize, chunk)

	return nil
}

func (w *World) CreateRegion(x, z int64) *Region {
	region := NewRegion(x, z)

	w.regions[fmt.Sprintf("%d:%d", x, z)] = region

	return region
}

func (w *World) Save() error {
	meta, err := nbt.Marshal(w.meta)

	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path.Join(w.cwd, "level.dat"), meta, 0777); err != nil {
		return err
	}

	for _, region := range w.regions {
		if !region.IsDirty() {
			continue
		}

		// TODO world save

		region.SetDirty(false)
	}

	return nil
}

func (w *World) Close() error {
	if err := w.Save(); err != nil {
		return err
	}

	return w.generator.Close()
}
