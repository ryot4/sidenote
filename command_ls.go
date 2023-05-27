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

func (c *LsCommand) Description() string {
	return "List notes"
}

func (c *LsCommand) setup(args []string, options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		output := c.flag.Output()
		fmt.Fprintf(output, "Usage: %s %s [-l] [-r] [-t] [name]\n", os.Args[0], c.Name())
		fmt.Fprintf(output, "\n%s.\n", c.Description())
		fmt.Fprintln(output, "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.BoolVar(&c.longFormat, "l", false, "Print modification time of entries")
	c.flag.BoolVar(&c.recurse, "r", false, "List directories recursively")
	c.flag.BoolVar(&c.sortByMtime, "t", false, "Sort entries by modification time (implies -l)")
	c.flag.Parse(args)

	if c.sortByMtime {
		c.longFormat = true
	}
}

func (c *LsCommand) Run(args []string, options *Options) error {
	c.setup(args, options)

	dir, err := options.CheckDirectory()
	if err != nil {
		return err
	}

	var name string
	switch c.flag.NArg() {
	case 0:
		name = ""
	case 1:
		name = c.flag.Arg(0)
	default:
		return ErrTooManyArgs
	}

	return c.list(dir, name)
}

func (c *LsCommand) list(dir *Directory, name string) error {
	path, err := dir.JoinPath(name)
	if err != nil {
		return err
	}
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		c.listDir(dir, name)
	} else {
		c.printFile(fi)
	}
	return nil
}

func (c *LsCommand) listDir(dir *Directory, name string) {
	items, err := dir.Readdir(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read %s (%s)\n", name, err)
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
				itemName := filepath.Join(name, fi.Name())
				fmt.Printf("\n%s:\n", itemName)
				c.listDir(dir, itemName)
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
