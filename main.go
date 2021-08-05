package main

import (
	"fmt"
	"io/fs"
	"os"
	"runtime/pprof"

	"github.com/jfontan/gofi/find"
	flag "github.com/spf13/pflag"
)

type Options struct {
	Workers        int
	Hidden         bool
	Regexp         bool
	MatchExtension string
	Profile        bool
}

func (o *Options) RegisterFlags() {
	flag.IntVarP(&o.Workers, "threads", "j", 0, "number of threads")
	flag.BoolVarP(&o.Hidden, "hidden", "H", false, "do not filter hidden files and directories")
	flag.BoolVarP(&o.Regexp, "regexp", "r", false, "search string is a regexp")
	flag.StringVarP(&o.MatchExtension, "extension", "e", "", "search for an extension")
	flag.BoolVar(&o.Profile, "profile", false, "write cpu profile to cpu.prof")
}

func (o *Options) ParseFlags() {
	flag.Parse()
}

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: gofi [options] [search] [path]")
		flag.PrintDefaults()
	}

	var opts Options
	opts.RegisterFlags()
	opts.ParseFlags()

	search := ""
	path := "."
	switch len(flag.Args()) {
	case 0:
	case 1:
		search = flag.Arg(0)
	case 2:
		search = flag.Arg(0)
		path = flag.Arg(1)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	fopts := find.Options{
		Hidden:         opts.Hidden,
		MatchExtension: opts.MatchExtension,
		Workers:        opts.Workers,
		Callback: func(path string, entry fs.DirEntry) error {
			fmt.Println(path)
			return nil
		},
	}

	if opts.Regexp {
		fopts.MatchRegexp = search
	} else {
		fopts.MatchString = search
	}

	f := find.New(path, fopts)

	if opts.Profile {
		stop, err := startPprof()
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			os.Exit(1)
		}
		defer stop()
	}

	_, err := f.Find()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}

func startPprof() (func(), error) {
	f, err := os.Create("cpu.prof")
	if err != nil {
		return nil, fmt.Errorf("could not create cpu profile: %w", err)
	}

	err = pprof.StartCPUProfile(f)
	if err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("could not start cpu profile: %w", err)
	}

	return func() {
		pprof.StopCPUProfile()
		f.Close()
	}, nil
}
