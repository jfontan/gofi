// package find implements a filesystem tree walker and filter. Similar to
// unix command "find".
//
// Using default options it returns all the files and directories in the
// provided path except hidden ones. If Option.Workers is 1 it sequentially
// reads all directories. If Option.Workers is greater than 0 it parallelizes
// the work with that amount of goroutines. Option.Workers = 0 means the same
// number as CPU count.
//
// NOTE: the list of files is not ordered and it can be different in every run.
//
// There are a number of options that can be used to filter the list of files.
// For example, to return only non hidden files and directories with extension
// ".go" use:
//
//		f := find.New("/path/to/files", find.Options{
//			Hidden: false,
//			MatchExtension: "go",
//		})
//		files, _ := f.Find()
//
// You can also specify a callback function that will be executed for all files
// and directories that match the filters.
package find

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/jfontan/fifo"
)

// Options change the default behavior of Find.
type Options struct {
	// Hidden tells if hidden files and directories are returned.
	Hidden bool
	// MatchString filters files that not contain the provided string.
	MatchString string
	// MathRegexp filters files that not match the provided regexp.
	MatchRegexp string
	// MatchExtension filters files that do not have the provided extension.
	MatchExtension string
	// Callback is executed for each file not filtered.
	Callback func(path string, entry fs.DirEntry) error

	// Workers specify the number of parallel goroutines to find and filter
	// files. 0 means the same number as CPUs.
	Workers int
}

// Find lists and filters files and directories in a provided path.
type Find struct {
	root    string
	opts    Options
	regexp  *regexp.Regexp
	workers int
}

// New creates a new Find. You can use find.Options{} for default options.
func New(path string, opts Options) *Find {
	workers := opts.Workers
	if workers == 0 {
		workers = runtime.NumCPU()
	}

	if opts.MatchExtension != "" &&
		!strings.HasPrefix(opts.MatchExtension, ".") {
		opts.MatchExtension = "." + opts.MatchExtension
	}

	return &Find{
		root:    path,
		workers: workers,
		opts:    opts,
	}
}

type result struct {
	dirs  []string
	files []string
	err   error
}

// Find searches the filesystem for files with the provided options.
func (f *Find) Find() ([]string, error) {
	if f.opts.MatchRegexp != "" {
		rg, err := regexp.Compile(f.opts.MatchRegexp)
		if err != nil {
			return nil, fmt.Errorf("invalid regexp: %w", err)
		}

		f.regexp = rg
	}

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

	paths := fifo.NewString()
	active := 0

	var files []string
	path := f.root
	value := true
	for {
		if value {
			select {
			case work <- path:
				path, value = paths.Pop()
				active++
			case res := <-result:
				active--
				if res.err == nil {
					paths.Push(res.dirs...)
					files = append(files, res.files...)
				} else {
					fmt.Printf("ERROR: %s\n", res.err.Error())
				}
			}
		} else {
			path, value = paths.Pop()
			if value {
				continue
			}

			if active == 0 {
				break
			}

			res := <-result
			active--
			if res.err == nil {
				paths.Push(res.dirs...)
				files = append(files, res.files...)
			} else {
				fmt.Printf("ERROR: %s\n", res.err.Error())
			}

			path, value = paths.Pop()
		}
	}

	close(work)
	return files, nil
}

func (f *Find) findSequential() ([]string, error) {
	paths := fifo.NewString()
	paths.Push(f.root)

	var files []string
	for !paths.Empty() {
		path, _ := paths.Pop()
		d, fs, err := f.process(path)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}
		files = append(files, fs...)
		paths.Push(d...)
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
		n := file.Name()
		if n == "." || n == ".." {
			continue
		}

		if !f.opts.Hidden && strings.HasPrefix(n, ".") {
			continue
		}

		fp := filepath.Join(path, n)
		if file.Type().IsDir() {
			dirs = append(dirs, fp)
		}

		if f.opts.MatchString != "" &&
			!strings.Contains(fp, f.opts.MatchString) {
			continue
		}

		if f.regexp != nil && !f.regexp.MatchString(fp) {
			continue
		}

		if f.opts.MatchExtension != "" &&
			!strings.HasSuffix(fp, f.opts.MatchExtension) {
			continue
		}

		if f.opts.Callback != nil {
			err := f.opts.Callback(fp, file)
			if err != nil {
				return nil, nil, fmt.Errorf("callback error: %w", err)
			}
		}

		files = append(files, fp)
	}

	return dirs, files, nil
}
