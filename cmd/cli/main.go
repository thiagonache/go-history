package main

import (
	"log"

	history "github.com/thiagonache/go-history"
)

func main() {
	r, err := history.NewRecorder()
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	r.SetPath("/tmp/history.log")
	go r.Session()
	select {
	case <-r.Context.Done():
		r.Shutdown()
	}
}
