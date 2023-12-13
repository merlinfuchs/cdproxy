package files

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/merlinfuchs/cdproxy/internal/config"
	"github.com/merlinfuchs/cdproxy/internal/db"
	"gopkg.in/guregu/null.v4"
)

type FileQueueEntry struct {
	ID                string
	OriginalURL       string
	OriginalExpiresAt null.Time
	ExpiresAt         null.Time
	Size              null.Int
	MaxSize           null.Int
	Metadata          map[string]string

	Done chan struct{}
}

func (fm *FileManager) queueWorker() {
	for {
		select {
		case <-fm.stop:
			return
		case entry := <-fm.queue:
			err := fm.processQueueEntry(entry)
			if err != nil {
				fmt.Println(err)
			}

			close(entry.Done)
			fm.Lock()
			delete(fm.doneChans, entry.ID)
			fm.Unlock()
		}
	}
}

func (fm *FileManager) processQueueEntry(entry *FileQueueEntry) error {
	if entry.OriginalExpiresAt.Valid && entry.OriginalExpiresAt.Time.Before(time.Now().UTC()) {
		return fmt.Errorf("original file expired before being processed")
	}

	if entry.ExpiresAt.Valid && entry.ExpiresAt.Time.Before(time.Now().UTC()) {
		return nil
	}

	size := entry.Size.Int64
	maxSize := entry.MaxSize.Int64
	var hash sql.NullString
	var contentType sql.NullString

	if maxSize == 0 {
		maxSize = int64(config.C.DefaultMaxSize)
	}

	if size <= maxSize {
		downloadRes, err := downloadFile(entry.OriginalURL)
		if err != nil {
			return fmt.Errorf("failed to download file: %w", err)
		}

		if len(downloadRes.Body) <= int(maxSize) {
			hashStr, err := hashFile(downloadRes.Body)
			if err != nil {
				return fmt.Errorf("failed to store file: %w", err)
			}

			err = fm.sftp.WriteFile(hashStr, downloadRes.Body)

			size = int64(len(downloadRes.Body))
			hash = sql.NullString{
				String: hashStr,
				Valid:  true,
			}
			contentType = sql.NullString{
				String: downloadRes.ContentType,
				Valid:  true,
			}
		}
	}

	err := fm.db.InsertFile(db.File{
		ID:                entry.ID,
		CreatedAt:         time.Now().UTC().Unix(),
		OriginalURL:       entry.OriginalURL,
		OriginalExpiresAt: null.NewInt(entry.OriginalExpiresAt.Time.Unix(), entry.OriginalExpiresAt.Valid),
		ExpiresAt:         null.NewInt(entry.ExpiresAt.Time.Unix(), entry.ExpiresAt.Valid),
		Size:              int(size),
		Hash:              hash,
		ContentType:       contentType,
		Metadata:          EncodeMetadata(entry.Metadata),
	})
	if err != nil {
		return fmt.Errorf("failed to insert file: %w", err)
	}

	return nil
}
