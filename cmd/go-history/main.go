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
		log.Fatal(err)
	}
	filenamePtr := flag.String("filename", defaultFilename.Name(), "filename to save recorded data")
	flag.Parse()

	fmt.Println("Welcome to history")
	fmt.Printf("See %s for recorded data\n", *filenamePtr)

	f, err := os.Create(*filenamePtr)
	if err != nil {
		log.Fatal(err)
	}
	HandleSigTerm(*filenamePtr)
	for {
		fmt.Print("$ ")
		err = history.Run(os.Stdin, f)
		errReceived := err != nil
		if errReceived && err.Error() == "exit" {
			break
		}
		if err != nil {
			log.Fatalf("%e %T", err, err)
		}
	}
}

// HandleSigTerm just avoid the program to crash
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
