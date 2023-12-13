package server

import (
	"fmt"
	"net/http"

	"github.com/merlinfuchs/cdproxy/internal/config"
	"github.com/merlinfuchs/cdproxy/internal/files"
)

type Server struct {
	fileManager *files.FileManager
}

func NewServer(fileManager *files.FileManager) Server {
	s := Server{
		fileManager: fileManager,
	}

	http.HandleFunc("/upload", s.handleUploadFile)
	http.HandleFunc("/download/", s.handleDownloadFile)

	return s
}

func (s *Server) Run() error {
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.C.Host, config.C.Port), nil)

	return err
}
