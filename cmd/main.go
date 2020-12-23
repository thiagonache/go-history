package main

import (
	"context"
	"history"
	"log"
	"os/signal"
)

func main() {
	r, err := history.NewRecorder()
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	r.Ctx, r.Stop = signal.NotifyContext(context.Background(), r.Signals...)
	go r.Session()
	select {
	case <-r.Ctx.Done():
		r.Shutdown()
	}
}
