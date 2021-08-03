package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	err := find()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}

var finder struct {
	paths []string
	pos   int
}

func find() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("wrong number of arguments")
	}

	finder.paths = []string{os.Args[1]}
	finder.pos = 0

	for finder.pos < len(finder.paths) {
		path := finder.paths[finder.pos]
		println(path, finder.pos, len(finder.paths))
		err := process(finder.paths[finder.pos])
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}
		finder.pos++
	}

	return nil
}

func process(path string) error {
	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range dir {
		if f.Name() == "." || f.Name() == ".." {
			continue
		}

		if !f.Type().IsDir() {
			continue
		}

		fp := filepath.Join(path, f.Name())
		// println(fp)
		finder.paths = append(finder.paths, fp)
	}

	return nil
}
