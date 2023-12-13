package pkg

import (
	"gopkg.in/guregu/null.v4"
)

// This only submit it to a queue, it will be processed asynchronously
// If the file is requested before it has been processed, cdproxy will block until it is ready
type SubmitFileRequestWire struct {
	OriginalURL       string            `json:"original_url"`
	OriginalExpiresAt null.Time         `json:"original_expires_at"` // when provided, cdproxy will not try to download the file after this date
	ExpiresAt         null.Time         `json:"expires_at"`
	Size              null.Int          `json:"size"`     // in bytes, when provided and bigger than max_size, cdproxy will not even attempt to download the file
	MaxSize           null.Int          `json:"max_size"` // in bytes
	Metadata          map[string]string `json:"metadata"`
	Wait              bool              `json:"wait"` // when true, cdproxy will block until the file has been processed
}

type SubmitFileResponseWire struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type FileInfoWire struct {
	ID                string            `json:"id"`
	URL               string            `json:"url"`
	OriginalURL       string            `json:"original_url"`
	OriginalExpiresAt null.Time         `json:"original_expires_at"`
	ExpiresAt         null.Time         `json:"expires_at"`
	Size              int               `json:"size"`
	Metadata          map[string]string `json:"metadata"`
}
