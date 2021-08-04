package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jfontan/gofind/find"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("wrong number of arguments")
		os.Exit(1)
	}

	sequential := find.New(os.Args[1], 1)

	start := time.Now()
	files, err := sequential.Find()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
	println("Time:", time.Since(start).String(), "Files:", len(files))

	parallel := find.New(os.Args[1], 8)

	start = time.Now()
	files, err = parallel.Find()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
	println("Time:", time.Since(start).String(), "Files:", len(files))
}
