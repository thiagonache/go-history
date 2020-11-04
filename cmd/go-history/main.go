package main

import (
	"flag"
	"fmt"
	"history"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defaultFilename, err := ioutil.TempFile("", "go-history-output-*.txt")
	if err != nil {
		log.Fatal(err)
	}
	filenamePtr := flag.String("filename", defaultFilename.Name(), "filename to save recorded data")
	flag.Parse()

	HandleSigTerm(*filenamePtr)

	fmt.Println("Welcome to history")
	f, err := os.Create(*filenamePtr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("See %s for recorded data\n", *filenamePtr)
	for {
		fmt.Print("$ ")
		err = history.Run(os.Stdin, f)
		// io.EOF means we should exit gracefully since we have nothing else to
		// read. It would happen if ctrl+d is pressed while reading the stdin or
		// if exit or quit commands are entered.
		if err == io.EOF {
			fmt.Printf("See %s for recorded data\n", *filenamePtr)
			os.Exit(0)
		}
		if err != nil {
			log.Fatalf("unexpected error: %v", err)
		}
	}
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
