package main

import "fmt"

type Parallel struct {
	root  string
	paths []string
	pos   int

	workers int
	active  int
}

func NewParallel(path string, workers int) *Parallel {
	return &Parallel{
		root:    path,
		workers: workers,
	}
}

type result struct {
	dirs  []string
	files []string
	err   error
}

func (p *Parallel) Find() ([]string, error) {
	work := make(chan string)
	result := make(chan result)
	for i := 0; i < p.workers; i++ {
		go p.worker(work, result)
	}

	p.paths = []string{p.root}
	p.pos = 0

	var files []string
	for {
		if p.pos < len(p.paths) {
			select {
			case work <- p.paths[p.pos]:
				p.pos++
				p.active++
			case res := <-result:
				p.active--
				if res.err == nil {
					p.paths = append(p.paths, res.dirs...)
					files = append(files, res.files...)
				} else {
					fmt.Printf("ERROR: %s\n", res.err.Error())
				}
			}
		} else {
			if p.active == 0 {
				break
			}

			res := <-result
			p.active--
			if res.err == nil {
				p.paths = append(p.paths, res.dirs...)
				files = append(files, res.files...)
			} else {
				fmt.Printf("ERROR: %s\n", res.err.Error())
			}

		}
	}

	close(work)
	return files, nil
}

func (p *Parallel) worker(work chan string, res chan result) {
	for w := range work {
		dirs, files, err := process(w)
		r := result{
			dirs:  dirs,
			files: files,
			err:   err,
		}
		res <- r
	}
}
