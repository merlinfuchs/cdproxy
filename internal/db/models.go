package db

import (
	"database/sql"

	"gopkg.in/guregu/null.v4"
)

type File struct {
	ID                string
	CreatedAt         int64
	OriginalURL       string
	OriginalExpiresAt null.Int
	ExpiresAt         null.Int
	Size              int
	Hash              sql.NullString
	ContentType       sql.NullString
	Metadata          string
}
