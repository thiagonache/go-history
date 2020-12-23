package main

import (
	"history"
	"log"
)

func main() {
	r, err := history.NewRecorder()
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	//r.ListenSignals(os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	r.ListenSignals()
	go r.Session()
	select {
	case <-r.Ctx.Done():
		r.Shutdown()
	}
}
