package world

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
)

type World struct {
	name           string
	folder         string
	store          WorldStore
	generator      WorldGenerator
	modifiedChunks map[string]*Chunk
}

func NewWorld(name string, folder string, store WorldStore, generator WorldGenerator) *World {
	return &World{
		name:           name,
		folder:         folder,
		store:          store,
		generator:      generator,
		modifiedChunks: make(map[string]*Chunk),
	}
}

func (w World) Name() string {
	return w.name
}

func (w *World) Initialize(storeOptions map[string]interface{}) error {
	worldFolder := path.Join(w.folder, w.name)

	if err := os.Mkdir(worldFolder, 0777); err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}

	w.folder = worldFolder

	if err := w.store.Initialize(worldFolder, storeOptions); err != nil {
		return err
	}

	return nil
}

func (w *World) GetChunk(x, z int64) (*Chunk, error) {
	chunk, ok := w.modifiedChunks[fmt.Sprintf("%d:%d", x, z)]

	if ok {
		return chunk, nil
	}

	chunk, ok, err := w.store.GetChunk(x, z)

	if err != nil {
		return nil, err
	}

	if !ok {
		newChunk, err := w.generator.GenerateChunk(x, z)

		if err != nil {
			return nil, err
		}

		w.modifiedChunks[fmt.Sprintf("%d:%d", x, z)] = newChunk

		chunk = newChunk
	}

	return chunk, nil
}

func (w *World) GetBlock(x, y, z int64) (*PaletteItem, error) {
	chunk, err := w.GetChunk(x/16, z/16)

	if err != nil {
		return nil, err
	}

	return chunk.GetBlock(x%16, y, z%16), nil
}

func (w *World) GenerateChunk(x, z int64) error {
	chunk, err := w.generator.GenerateChunk(x, z)

	if err != nil {
		return err
	}

	w.modifiedChunks[fmt.Sprintf("%d:%d", x, z)] = chunk

	return nil
}

func (w *World) Save() error {
	for _, chunk := range w.modifiedChunks {
		if err := w.store.SetChunk(chunk.x, chunk.z, *chunk); err != nil {
			return err
		}
	}

	return nil
}

func (w *World) Close() error {
	if err := w.Save(); err != nil {
		return err
	}

	if err := w.generator.Close(); err != nil {
		return err
	}

	return w.store.Close()
}
