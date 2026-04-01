package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New creates a new storage instance with given path to sqlite file.
// If file doesn't exist, it will be created.
// Creates pack_sizes table and seeds default values if table is empty.
func New(storagePath string, defaultSizes []int) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Create table if it doesn't exist.
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS pack_sizes (
		size INTEGER PRIMARY KEY CHECK(size > 0)
	)`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Seed default pack sizes if table is empty.
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM pack_sizes").Scan(&count)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if count == 0 {
		tx, err := db.Begin()
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		for _, size := range defaultSizes {
			_, err = tx.Exec("INSERT INTO pack_sizes (size) VALUES (?)", size)
			if err != nil {
				tx.Rollback()
				db.Close()
				return nil, fmt.Errorf("%s: %w", op, err)
			}
		}

		if err := tx.Commit(); err != nil {
			db.Close()
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &Storage{db: db}, nil
}

// GetPackSizes returns available pack sizes sorted in ascending order.
func (s *Storage) GetPackSizes() ([]int, error) {
	rows, err := s.db.Query("SELECT size FROM pack_sizes ORDER BY size")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sizes []int
	for rows.Next() {
		var size int
		if err := rows.Scan(&size); err != nil {
			return nil, err
		}
		sizes = append(sizes, size)
	}
	return sizes, rows.Err()
}

// UpdatePackSizes replaces all pack sizes with the given list.
// Runs in a transaction — either all sizes are replaced, or none.
func (s *Storage) UpdatePackSizes(sizes []int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM pack_sizes")
	if err != nil {
		return err
	}

	for _, size := range sizes {
		_, err = tx.Exec("INSERT INTO pack_sizes (size) VALUES (?)", size)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Close closes the database connection.
func (s *Storage) Close() error {
	return s.db.Close()
}