package world

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/Tnze/go-mc/nbt"
	log "github.com/anchormc/anchor/src/api/logger"
	"github.com/anchormc/anchor/src/api/world"
	"github.com/anchormc/anchor/src/game/generators"
)

type World struct {
	cwd       string
	meta      WorldMeta
	generator world.WorldGenerator
	regions   map[string]world.Region
}

type WorldOptions struct {
	Folder    string
	Generator world.WorldGenerator
	Height    int
}

type WorldMeta struct {
	Name      string `nbt:"name"`
	Generator string `nbt:"generator"`
	Height    int32  `nbt:"height"`
}

func LoadWorldFromFolder(cwd string) (world.World, error) {
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
		regions:   make(map[string]world.Region),
	}, nil
}

func NewWorld(name string, options WorldOptions) world.World {
	return &World{
		cwd: options.Folder,
		meta: WorldMeta{
			Name:      name,
			Generator: "flat",
			Height:    256,
		},
		generator: options.Generator,
		regions:   make(map[string]world.Region),
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

func (w *World) GetRegion(x, z int64) (world.Region, bool, error) {
	region, ok := w.regions[fmt.Sprintf("%d:%d", x, z)]

	if ok {
		return region, true, nil
	}

	_, err := os.Stat(path.Join(w.cwd, "region", fmt.Sprintf("region.%d.%d.bin", x, z)))

	if err != nil {
		return nil, false, nil
	}

	f, err := os.Open(path.Join(w.cwd, "region", fmt.Sprintf("region.%d.%d.bin", x, z)))

	if err != nil {
		return nil, false, err
	}

	defer f.Close()

	newRegion := &Region{}

	if err = newRegion.Decode(f); err != nil {
		return nil, true, err
	}

	w.regions[fmt.Sprintf("%d:%d", x, z)] = newRegion

	return newRegion, true, nil
}

func (w *World) GetChunk(x, z int64) (world.Chunk, bool, error) {
	region, ok, err := w.GetRegion(x/world.RegionSize, z/world.RegionSize)

	if err != nil || !ok {
		return nil, ok, err
	}

	return region.GetChunk(x%world.RegionSize, z%world.RegionSize), true, nil
}

func (w *World) GetChunkOrGenerate(x, z int64) (world.Chunk, error) {
	region, ok, err := w.GetRegion(x/world.RegionSize, z/world.RegionSize)

	if err != nil {
		return nil, err
	}

	if !ok {
		region = w.CreateRegion(x, z)

		if err := w.GenerateChunk(x, z); err != nil {
			return nil, err
		}
	}

	return region.GetChunk(x%world.RegionSize, z%world.RegionSize), nil
}

func (w *World) GetBlock(x, y, z int64) (*world.Block, bool, error) {
	c, ok, err := w.GetChunk(x/int64(world.ChunkSize), z/int64(world.ChunkSize))

	if err != nil || !ok {
		return nil, ok, err
	}

	return c.GetBlock(x%int64(world.ChunkSize), y, z%int64(world.ChunkSize)), true, nil
}

func (w *World) GenerateChunk(x, z int64) error {
	region, ok, err := w.GetRegion(x/world.RegionSize, z/world.RegionSize)

	if err != nil {
		return err
	}

	if !ok {
		region = w.CreateRegion(x/world.RegionSize, z/world.RegionSize)
	}

	c := region.GetChunk(x, z)

	if c == nil {
		c = region.CreateChunk(x, z)
	}

	if err := w.generator.GenerateChunk(c); err != nil {
		return err
	}

	return nil
}

func (w *World) CreateRegion(x, z int64) world.Region {
	region := NewRegion(x, z)

	w.regions[fmt.Sprintf("%d:%d", x, z)] = region

	return region
}

func (w *World) Save() error {
	meta, err := nbt.Marshal(w.meta)

	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(path.Join(w.cwd, "level.dat"), meta, 0777); err != nil {
		return err
	}

	for _, region := range w.regions {
		log.Debugf("world", "Saving region (%d, %d) for world '%s'\n", region.X(), region.Z(), w.Name())

		f, err := os.OpenFile(path.Join(w.cwd, "region", fmt.Sprintf("region.%d.%d.dat", region.X(), region.Z())), os.O_CREATE|os.O_RDWR, 0777)

		if err != nil {
			return err
		}

		if err = region.Encode(f); err != nil {
			return err
		}

		if err = f.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (w *World) Close() error {
	if err := w.Save(); err != nil {
		return err
	}

	return w.generator.Close()
}

var _ world.World = &World{}
