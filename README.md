# gofi

This project implements a filesystem tree walker and filter. Similar to unix command "find". Can be used as a library for go projects and has a simple command line.

It was developed mainly to experiment with parallel tree walking and command line interfaces.

## Command

The command functionality is very limited and inspired by [fd](https://github.com/sharkdp/fd). Use "fd" instead of this project for a fully featured parallel find replacement.

```
Usage: gofi [options] [search] [path]
  -e, --extension string   search for an extension
  -H, --hidden             do not filter hidden files and directories
      --profile            write cpu profile to cpu.prof
  -r, --regexp             search string is a regexp
  -j, --threads int        number of threads
```

To search for all files in another directory use an empty search string `""`. For example:

```
$ gofi "" /another/directory
```

To install it you can download the binaries from the latest release: https://github.com/jfontan/gofi/releases/latest

You can also install it with go:

```
$ go install github.com/jfontan/gofi
```

## Library

Documentation: https://pkg.go.dev/github.com/jfontan/gofi/find

Using default options it returns all the files and directories in the provided path except hidden ones. If `Options.Workers` is `1` it sequentially reads all directories. If `Options.Workers` is greater than `0` it parallelizes the work with that amount of goroutines. `Options.Workers` = `0` means the same number as CPU count.

NOTE: the list of files is not ordered and it can be different in every run.

There are a number of options that can be used to filter the list of files.  For example, to return only non hidden files and directories with extension ".go" use:

```go
import "github.com/jfontan/gofi/find"

...

    f := find.New("/path/to/files", find.Options{
        Hidden: false,
        MatchExtension: "go",
    })
    files, _ := f.Find()
```

You can also specify a callback function that will be executed for all files and directories that match the filters.