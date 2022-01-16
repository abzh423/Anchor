package store

import (
	"bytes"
	"database/sql"
	"fmt"
	"path"

	_ "github.com/mattn/go-sqlite3"

	"github.com/golangminecraft/minecraft-server/src/api/world"
)

const (
	createChunksTableStatement = `CREATE TABLE IF NOT EXISTS "chunks" (
		"x" INTEGER NOT NULL,
		"z" INTEGER NOT NULL,
		"data" BLOB NOT NULL,
		PRIMARY KEY("x", "z")
	)`
)

type SQLiteStore struct {
	conn *sql.DB
}

func NewSQLiteStore() world.WorldStore {
	return &SQLiteStore{
		conn: nil,
	}
}

func (s *SQLiteStore) Initialize(folder string, _ map[string]interface{}) error {
	conn, err := sql.Open("sqlite3", path.Join(folder, "world.db"))

	if err != nil {
		return err
	}

	if _, err := conn.Exec(createChunksTableStatement); err != nil {
		conn.Close()

		return err
	}

	s.conn = conn

	return nil
}

func (s *SQLiteStore) GetChunk(x, z int64) (*world.Chunk, bool, error) {
	rows, err := s.conn.Query("SELECT data FROM chunks WHERE x = ? AND z = ?", x, z)

	if err != nil {
		return nil, false, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, false, nil
	}

	var data []byte

	if err = rows.Scan(&data); err != nil {
		return nil, false, err
	}

	if err = rows.Err(); err != nil {
		return nil, false, err
	}

	chunk := world.NewChunk(x, z)

	if err := chunk.Decode(bytes.NewReader(data)); err != nil {
		return nil, false, err
	}

	return chunk, true, nil
}

func (s *SQLiteStore) SetChunk(x, z int64, chunk world.Chunk) error {
	buf := &bytes.Buffer{}

	if err := chunk.Encode(buf); err != nil {
		return err
	}

	rows, err := s.conn.Query("SELECT COUNT(*) FROM chunks WHERE x = ? AND z = ?", x, z)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf("unexpected at least 1 row, got none")
	}

	var count int

	if err := rows.Scan(&count); err != nil {
		return err
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if err := rows.Close(); err != nil {
		return err
	}

	if count >= 1 {
		_, err = s.conn.Exec("UPDATE chunks SET data = ? WHERE x = ? AND z = ?", buf.Bytes(), x, z)

		return err
	}

	_, err = s.conn.Exec("INSERT INTO chunks (x, z, data) VALUES (?, ?, ?)", x, z, buf.Bytes())

	return err
}

func (s *SQLiteStore) Close() error {
	return s.conn.Close()
}
