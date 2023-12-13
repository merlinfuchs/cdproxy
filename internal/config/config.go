package config

import (
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)

var C = Config{
	Host:                   "127.0.0.1",
	Port:                   8080,
	PublicURL:              "http://localhost:8080",
	DBFileName:             "cdproxy.db",
	MaxQueueSize:           100,
	NumWorkers:             runtime.NumCPU(),
	BrotliCompressionLevel: 7,
	FilePath:               "./files",
	DownloadTimeout:        30,
	DefaultMaxSize:         100 * 1024 * 1024, // 100 MB
}

func InitConifg() error {
	fileContent, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(fileContent, &C)
	if err != nil {
		return err
	}
	return nil
}
