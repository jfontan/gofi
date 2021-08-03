package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Sequential struct {
	root  string
	paths []string
	pos   int
}

func NewSequential(path string) *Sequential {
	return &Sequential{
		root: path,
	}
}

func (s *Sequential) Find() ([]string, error) {
	s.paths = []string{s.root}
	s.pos = 0

	var files []string
	for s.pos < len(s.paths) {
		d, f, err := process(s.paths[s.pos])
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}
		files = append(files, f...)
		s.paths = append(s.paths, d...)
		s.pos++
	}

	return files, nil
}

func process(path string) ([]string, []string, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, nil, err
	}

	var files []string
	var dirs []string
	for _, f := range dir {
		if f.Name() == "." || f.Name() == ".." {
			continue
		}

		fp := filepath.Join(path, f.Name())
		if !f.Type().IsDir() {
			files = append(files, fp)
			continue
		}

		dirs = append(dirs, fp)
	}

	return dirs, files, nil
}
