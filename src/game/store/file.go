package store

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"

	"github.com/golangminecraft/minecraft-server/src/api/world"
)

type FileStore struct {
	folder string
}

func NewFileStore() world.WorldStore {
	return &FileStore{
		folder: "",
	}
}

func (s *FileStore) Initialize(folder string, _ map[string]interface{}) error {
	s.folder = folder

	if err := os.Mkdir(path.Join(folder, "region"), 0777); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	return nil
}

func (s *FileStore) GetChunk(x, z int64) (*world.Chunk, bool, error) {
	data, err := ioutil.ReadFile(path.Join(s.folder, "region", fmt.Sprintf("%d-%d.bin", x, z)))

	if err != nil {
		return nil, false, err
	}

	chunk := world.NewChunk(x, z)

	if err := chunk.Decode(bytes.NewReader(data)); err != nil {
		return nil, false, err
	}

	return chunk, true, nil
}

func (s *FileStore) SetChunk(x, z int64, chunk world.Chunk) error {
	buf := &bytes.Buffer{}

	if err := chunk.Encode(buf); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path.Join(s.folder, "region", fmt.Sprintf("chunk.%d.%d.bin", x, z)), buf.Bytes(), 0777); err != nil {
		return err
	}

	return nil
}

func (s *FileStore) Close() error {
	return nil
}
