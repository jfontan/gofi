# gofi

Package find implements a filesystem tree walker and filter. Similar to unix command "find".

Using default options it returns all the files and directories in the provided path except hidden ones. If Option.Workers is 1 it sequentially reads all directories. If Option.Workers is greater than 0 it parallelizes the work with that amount of goroutines. Option.Workers = 0 means the same number as CPU count.

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