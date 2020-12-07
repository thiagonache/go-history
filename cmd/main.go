package main

import (
	"fmt"
	"history"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	HandleSigTerm()
	err := history.RecordSession(os.Stdin, os.Stdout, nil)
	// err.Is or err.As
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
}

// HandleSigTerm just avoid the program to crash by handling sigterm.
func HandleSigTerm() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Sigterm received. Gracefully shutting down")
		fmt.Printf("\rSee recorded data at %s\n", history.LogFile)
		os.Exit(0)
	}()
}
