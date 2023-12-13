package files

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/merlinfuchs/cdproxy/internal/config"
	"gopkg.in/guregu/null.v4"
)

type File struct {
	ID                string
	CreatedAt         time.Time
	OriginalURL       string
	OriginalExpiresAt null.Time
	ExpiresAt         null.Time
	Size              int
	Hash              null.String
	ContentType       null.String
	Metadata          map[string]string
}

func (f *File) Download() (downloadFileResult, error) {
	if f.Hash.Valid {
		body, err := readFile(f.Hash.String)
		if err != nil {
			return downloadFileResult{}, fmt.Errorf("failed to read file: %w", err)
		}

		return downloadFileResult{
			Body:        body,
			ContentType: f.ContentType.String,
		}, nil
	}

	res, err := downloadFile(f.OriginalURL)
	if err != nil {
		return downloadFileResult{}, fmt.Errorf("failed to download file: %w", err)
	}

	return res, nil
}

func downloadFile(originalURL string) (downloadFileResult, error) {
	res := downloadFileResult{}

	req, err := http.NewRequest("GET", originalURL, nil)
	if err != nil {
		return res, fmt.Errorf("failed to create request: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.C.DownloadTimeout)*time.Second)
	req = req.WithContext(ctx)
	defer cancel()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return res, fmt.Errorf("failed to download file: %w", err)
	}

	res.ContentType = resp.Header.Get("Content-Type")

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, fmt.Errorf("failed to read response body: %w", err)
	}

	res.Body = body
	return res, nil
}

func storeFile(body []byte) (string, error) {
	hasher := sha256.New()
	hasher.Write(body)
	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	path := filepath.Join(config.C.FilePath, hash)
	file, err := os.Create(path)

	writer := brotli.NewWriterLevel(file, config.C.BrotliCompressionLevel)
	_, err = writer.Write(body)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to brotli writer file: %w", err)
	}

	return hash, nil
}

func deleteFile(hash string) error {
	path := filepath.Join(config.C.FilePath, hash)
	err := os.Remove(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return nil
}

func readFile(hash string) ([]byte, error) {
	path := filepath.Join(config.C.FilePath, hash)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	reader := brotli.NewReader(file)
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return body, nil
}

type downloadFileResult struct {
	Body        []byte
	ContentType string
}
