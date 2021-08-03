package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("wrong number of arguments")
		os.Exit(1)
	}

	sequential := NewSequential(os.Args[1])
	files, err := sequential.Find()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
	println("Files:", len(files))
}
