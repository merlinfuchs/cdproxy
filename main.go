package main

import (
	"log"

	"github.com/merlinfuchs/cdproxy/internal/config"
	"github.com/merlinfuchs/cdproxy/internal/db"
	"github.com/merlinfuchs/cdproxy/internal/files"
	"github.com/merlinfuchs/cdproxy/internal/server"
)

func main() {
	err := config.InitConifg()
	if err != nil {
		panic(err)
	}

	db, err := db.New()
	if err != nil {
		panic(err)
	}

	err = db.InitDatabase()
	if err != nil {
		panic(err)
	}

	fileManager := files.NewFileManager(db)
	fileManager.StartWorkers()

	s := server.NewServer(fileManager)

	log.Printf("Starting server on http://%s:%d\n", config.C.Host, config.C.Port)

	err = s.Run()
	if err != nil {
		panic(err)
	}
}
