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

func (c *LsCommand) Run(args []string, options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.BoolVar(&c.longFormat, "l", false, "long format")
	c.flag.BoolVar(&c.recurse, "r", false, "list directories recursively")
	c.flag.BoolVar(&c.sortByMtime, "t", false, "sort by modification time")
	c.flag.Parse(args)

	dir, err := OpenDirectory(options.noteDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	listPaths := c.flag.Args()
	if len(listPaths) == 0 {
		listPaths = append(listPaths, "")
	}

	var files []os.FileInfo
	var dirPaths []string
	for _, path := range listPaths {
		fi, err := dir.Stat(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot stat %s: %s\n", path, err)
			continue
		}
		if fi.IsDir() {
			dirPaths = append(dirPaths, path)
		} else {
			files = append(files, fi)
		}
	}

	c.listItems(files)
	for i, path := range dirPaths {
		if len(listPaths) > 1 && i > 0 {
			fmt.Println()
		}
		if len(listPaths) > 1 {
			fmt.Printf("%s:\n", path)
		}
		c.listDir(dir, path)
	}
}

func (c *LsCommand) listDir(dir *Directory, path string) {
	children, err := dir.Readdir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read %s: %s\n", path, err)
	}
	c.listItems(children)
	if c.recurse {
		for _, item := range children {
			if item.IsDir() {
				itemPath := filepath.Join(path, item.Name())
				fmt.Printf("\n%s:\n", itemPath)
				c.listDir(dir, itemPath)
			}
		}
	}
}

func (c *LsCommand) listItems(items []os.FileInfo) {
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
