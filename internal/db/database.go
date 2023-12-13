package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/merlinfuchs/cdproxy/internal/config"
)

type Database struct {
	db *sql.DB
}

func New() (*Database, error) {
	db, err := sql.Open("sqlite3", config.C.DBFileName)
	if err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}

func (d *Database) InitDatabase() error {
	_, err := d.db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			id TEXT PRIMARY KEY,
			created_at INTEGER NOT NULL,
			original_url TEXT NOT NULL,
			original_expires_at INTEGER,
			expires_at INTEGER,
			size INTEGER NOT NULL,
			hash TEXT,
			content_type TEXT,
			metadata TEXT
		)
	`)
	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}
