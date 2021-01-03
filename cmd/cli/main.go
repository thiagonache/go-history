package main

import (
	"log"

	history "github.com/thiagonache/go-history"
)

func main() {
	r, err := history.NewRecorder(
		history.WithLogPath("/tmp/history.log"),
		history.WithLogPermission(0600),
	)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	go r.Session()
	r.WaitForExit()
}
