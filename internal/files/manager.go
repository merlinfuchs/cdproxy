package files

import (
	"fmt"
	"sync"
	"time"

	"github.com/merlinfuchs/cdproxy/internal/config"
	"github.com/merlinfuchs/cdproxy/internal/db"
	"github.com/merlinfuchs/cdproxy/internal/sftp"
	"github.com/merlinfuchs/cdproxy/internal/util"
	"gopkg.in/guregu/null.v4"
)

type FileManager struct {
	sync.Mutex

	db        *db.Database
	sftp      *sftp.SFTPClient
	queue     chan *FileQueueEntry
	stop      chan struct{}
	doneChans map[string]chan struct{}
}

func NewFileManager(db *db.Database, sftp *sftp.SFTPClient) *FileManager {
	return &FileManager{
		db:        db,
		sftp:      sftp,
		queue:     make(chan *FileQueueEntry, config.C.MaxQueueSize),
		stop:      make(chan struct{}),
		doneChans: make(map[string]chan struct{}),
	}
}

func (fm *FileManager) StartWorkers() {
	for i := 0; i < config.C.NumWorkers; i++ {
		go fm.queueWorker()
	}

	go fm.expiryWorker()
}

func (fm *FileManager) StopWorkers() {
	close(fm.stop)
}

type FileProcessRequest struct {
	OriginalURL       string
	OriginalExpiresAt null.Time
	ExpiresAt         null.Time         `json:"expires_at"`
	Size              null.Int          `json:"size"`
	MaxSize           null.Int          `json:"max_size"`
	Metadata          map[string]string `json:"metadata"`
}

func (fm *FileManager) ProcessFile(req FileProcessRequest) *FileQueueEntry {
	entry := &FileQueueEntry{
		ID:                util.UniqueID(),
		OriginalURL:       req.OriginalURL,
		OriginalExpiresAt: req.OriginalExpiresAt,
		ExpiresAt:         req.ExpiresAt,
		Size:              req.Size,
		MaxSize:           req.MaxSize,
		Metadata:          req.Metadata,
		Done:              make(chan struct{}),
	}

	fm.Lock()
	fm.doneChans[entry.ID] = entry.Done
	fm.Unlock()

	fm.queue <- entry
	return entry
}

func (fm *FileManager) WaitForFile(id string) {
	fm.Lock()
	c, ok := fm.doneChans[id]
	fm.Unlock()

	if ok {
		<-c
	}
}

func (fm *FileManager) GetFile(id string) (*File, error) {
	fm.WaitForFile(id)

	file, err := fm.db.GetFile(id)
	if err != nil {
		return nil, err
	}

	metadata, err := DecodeMetadata(file.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}

	return &File{
		ID:                file.ID,
		CreatedAt:         time.Unix(file.CreatedAt, 0),
		OriginalURL:       file.OriginalURL,
		OriginalExpiresAt: null.NewTime(time.Unix(file.OriginalExpiresAt.Int64, 0), file.OriginalExpiresAt.Valid),
		ExpiresAt:         null.NewTime(time.Unix(file.ExpiresAt.Int64, 0), file.ExpiresAt.Valid),
		Size:              file.Size,
		Hash:              null.String{NullString: file.Hash},
		ContentType:       null.String{NullString: file.ContentType},
		Metadata:          metadata,
	}, nil
}
