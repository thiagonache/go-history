package main

import (
	"history"
	"os"
)

func main() {
	history.Run(os.Stdout)
}
