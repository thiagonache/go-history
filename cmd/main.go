package main

import (
	"flag"
	"fmt"
	"history"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defaultFilename, err := ioutil.TempFile("", "go-history-output-*.txt")
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	filenamePtr := flag.String("filename", defaultFilename.Name(), "filename to save recorded data")
	flag.Parse()

	HandleSigTerm(*filenamePtr)

	fmt.Println("Welcome to history")
	f, err := os.Create(*filenamePtr)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("See %s for recorded data\n", *filenamePtr)

	err = history.RecordSession(os.Stdin, os.Stdout, f)
	// err.Is or err.As
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("See %s for recorded data\n", *filenamePtr)
}

// HandleSigTerm just avoid the program to crash by handling sigterm.
func HandleSigTerm(filename string) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Sigterm received. Gracefully shutting down")
		fmt.Printf("\rSee recorded data at %s\n", filename)
		os.Exit(0)
	}()
}
