package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/merlinfuchs/cdproxy/internal/config"
	"github.com/merlinfuchs/cdproxy/internal/files"
	"github.com/merlinfuchs/cdproxy/pkg"
)

func (s *Server) handleSubmitFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to read request body: %v\n", err)
		return
	}

	req := pkg.SubmitFileRequestWire{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("failed to unmarshal request body: %v\n", err)
		return
	}

	entry := s.fileManager.ProcessFile(files.FileProcessRequest{
		OriginalURL:       req.OriginalURL,
		OriginalExpiresAt: req.OriginalExpiresAt,
		ExpiresAt:         req.ExpiresAt,
		Size:              req.Size,
		MaxSize:           req.MaxSize,
		Metadata:          req.Metadata,
	})

	if req.Wait {
		<-entry.Done
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(pkg.SubmitFileResponseWire{
		ID:  entry.ID,
		URL: config.C.PublicURL + "/download/" + entry.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to encode response: %v\n", err)
		return
	}
}

func (s *Server) handleDownloadFile(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Path[len("/download/"):]

	file, err := s.fileManager.GetFile(fileID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to get file: %v\n", err)
		return
	}

	res, err := file.Download()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to download file: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", res.ContentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(res.Body)))
	w.Header().Set("Content-Disposition", "inline")
	w.Write(res.Body)
}
