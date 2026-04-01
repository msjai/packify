package sqlite

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

// New storage instance with given path to sqlite file.
// If file doesn't exist, it will be created.
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// GetPackSizes returns available pack sizes.
func (s Storage) GetPackSizes() ([]int, error) {
	// TODO implement me
	panic("implement me")
}

func (s Storage) UpdatePackSizes(sizes []int) error {
	// TODO implement me
	panic("implement me")
}
