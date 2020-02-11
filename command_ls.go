package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type LsCommand struct {
	flag *flag.FlagSet

	longFormat  bool
	recurse     bool
	sortByMtime bool
}

func (c *LsCommand) Name() string {
	return "ls"
}

func (c *LsCommand) setup(args []string, options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		fmt.Fprintf(c.flag.Output(), "Usage: %s [-l] [-r] [-t] [name]\n", c.Name())
		fmt.Fprintln(c.flag.Output(), "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.BoolVar(&c.longFormat, "l", false, "Print modification time of entries")
	c.flag.BoolVar(&c.recurse, "r", false, "List directories recursively")
	c.flag.BoolVar(&c.sortByMtime, "t", false, "Sort entries by modification time")
	c.flag.Parse(args)
}

func (c *LsCommand) Run(args []string, options *Options) {
	c.setup(args, options)

	dir, err := openDirectory(options.noteDir)
	if err != nil {
		exitWithError(err)
	}

	var listPath string
	if c.flag.NArg() > 1 {
		exitWithSyntaxError("too many arguments")
	} else if c.flag.NArg() == 1 {
		listPath = c.flag.Arg(0)
	} else {
		listPath = ""
	}

	err = c.list(dir, listPath)
	if err != nil {
		exitWithError(err)
	}
}

func (c *LsCommand) list(dir *Directory, path string) error {
	realPath, err := dir.FilePath(path)
	if err != nil {
		return err
	}
	fi, err := os.Stat(realPath)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		c.listDir(dir, path)
	} else {
		c.printFile(fi)
	}
	return nil
}

func (c *LsCommand) listDir(dir *Directory, path string) {
	items, err := dir.Readdir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read %s: %s\n", path, err)
		return
	}

	if c.sortByMtime {
		sort.Slice(items, func(i, j int) bool {
			return items[i].ModTime().After(items[j].ModTime())
		})
	} else {
		sort.Slice(items, func(i, j int) bool {
			return items[i].Name() < items[j].Name()
		})
	}

	for _, fi := range items {
		c.printFile(fi)
	}

	if c.recurse {
		for _, fi := range items {
			if fi.IsDir() {
				itemPath := filepath.Join(path, fi.Name())
				fmt.Printf("\n%s:\n", itemPath)
				c.listDir(dir, itemPath)
			}
		}
	}
}

func (c *LsCommand) printFile(fi os.FileInfo) {
	name := fi.Name()
	if fi.IsDir() {
		name += "/"
	}
	if c.longFormat {
		fmt.Printf("%s %s\n", formatTime(fi.ModTime()), name)
	} else {
		fmt.Println(name)
	}
}

func formatTime(t time.Time) string {
	if t.Year() == time.Now().Year() {
		return t.Format("Jan _2 15:04")
	} else {
		return t.Format("Jan _2  2006")
	}
}
