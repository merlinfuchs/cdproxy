package util

import (
	"crypto/rand"
	"fmt"
	"log"
)

func UniqueID() string {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		log.Fatalf("error while generating random string: %s", err)
	}

	return fmt.Sprintf("%x", buf)
}
