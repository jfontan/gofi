package find

import (
	"fmt"
	"os"
	"path/filepath"
)

type Find struct {
	root  string
	paths []string
	pos   int

	workers int
	active  int
}

func New(path string, workers int) *Find {
	return &Find{
		root:    path,
		workers: workers,
	}
}

type result struct {
	dirs  []string
	files []string
	err   error
}

func (f *Find) Find() ([]string, error) {
	if f.workers > 1 {
		return f.findParallel()
	}
	return f.findSequential()
}

func (f *Find) findParallel() ([]string, error) {
	work := make(chan string)
	result := make(chan result)
	for i := 0; i < f.workers; i++ {
		go f.worker(work, result)
	}

	f.paths = []string{f.root}
	f.pos = 0
	f.active = 0

	var files []string
	for {
		if f.pos < len(f.paths) {
			select {
			case work <- f.paths[f.pos]:
				f.pos++
				f.active++
			case res := <-result:
				f.active--
				if res.err == nil {
					f.paths = append(f.paths, res.dirs...)
					files = append(files, res.files...)
				} else {
					fmt.Printf("ERROR: %s\n", res.err.Error())
				}
			}
		} else {
			if f.active == 0 {
				break
			}

			res := <-result
			f.active--
			if res.err == nil {
				f.paths = append(f.paths, res.dirs...)
				files = append(files, res.files...)
			} else {
				fmt.Printf("ERROR: %s\n", res.err.Error())
			}

		}
	}

	close(work)
	return files, nil
}

func (f *Find) findSequential() ([]string, error) {
	f.paths = []string{f.root}
	f.pos = 0

	var files []string
	for f.pos < len(f.paths) {
		d, fs, err := f.process(f.paths[f.pos])
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}
		files = append(files, fs...)
		f.paths = append(f.paths, d...)
		f.pos++
	}

	return files, nil
}

func (f *Find) worker(work chan string, res chan result) {
	for w := range work {
		dirs, files, err := f.process(w)
		r := result{
			dirs:  dirs,
			files: files,
			err:   err,
		}
		res <- r
	}
}

func (f *Find) process(path string) ([]string, []string, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, nil, err
	}

	var files []string
	var dirs []string
	for _, file := range dir {
		if file.Name() == "." || file.Name() == ".." {
			continue
		}

		fp := filepath.Join(path, file.Name())
		if !file.Type().IsDir() {
			files = append(files, fp)
			continue
		}

		dirs = append(dirs, fp)
	}

	return dirs, files, nil
}
