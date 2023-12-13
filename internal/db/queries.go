package db

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

func (db *Database) InsertFile(f File) error {
	_, err := db.db.Exec(`
		INSERT INTO files (
			id,
			created_at,
			original_url,
			original_expires_at,
			expires_at,
			size,
			hash,
			content_type,
			metadata
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, f.ID, f.CreatedAt, f.OriginalURL, f.OriginalExpiresAt, f.ExpiresAt, f.Size, f.Hash, f.ContentType, f.Metadata)

	return err
}

func (db *Database) GetFile(id string) (File, error) {
	var res File

	err := db.db.
		QueryRow(`SELECT * FROM files WHERE id = ?`, id).
		Scan(&res.ID, &res.CreatedAt, &res.OriginalURL, &res.OriginalExpiresAt, &res.ExpiresAt, &res.Size, &res.Hash, &res.ContentType, &res.Metadata)

	if err != nil {
		return File{}, err
	}

	return res, nil
}

func (db *Database) GetExpiredFileHashes() ([]string, error) {
	rows, err := db.db.Query(`
		SELECT hash, MAX(expires_at) as max_expires_at FROM files 
		WHERE hash IS NOT NULL
		GROUP BY hash 
		HAVING max_expires_at <= ? AND SUM(expires_at IS NULL) = 0
	`, time.Now().UTC().Unix())
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := []string{}
	for rows.Next() {
		var hash string
		var maxExpiresAt string
		err := rows.Scan(&hash, &maxExpiresAt)
		if err != nil {
			return nil, err
		}

		res = append(res, hash)
	}

	return res, nil
}

func (db *Database) RemoveHashFromFiles(hash string) error {
	_, err := db.db.Exec(`UPDATE files SET hash = NULL WHERE hash = ?`, hash)
	return err
}
