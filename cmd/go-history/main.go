package main

import (
	"fmt"
	"history"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Welcome to history")
	HandleSigTerm()
	history.Run()
}

// HandleSigTerm just avoid the program to crash
func HandleSigTerm() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Sigterm received. Gracefully shutting down")
		os.Exit(0)
	}()
}
