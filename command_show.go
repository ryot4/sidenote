package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type ShowCommand struct {
	flag *flag.FlagSet
}

func (c *ShowCommand) Name() string {
	return "show"
}

func (c *ShowCommand) setup(args []string, _options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		fmt.Fprintf(c.flag.Output(), "Usage: %s path-to-file\n", c.Name())
	}
	c.flag.Parse(args)
}

func (c *ShowCommand) Run(args []string, options *Options) {
	c.setup(args, options)

	dir, err := OpenDirectory(options.noteDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var filePath string
	if c.flag.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "too many arguments")
		os.Exit(2)
	} else if c.flag.NArg() == 1 {
		filePath = c.flag.Arg(0)
	} else {
		fmt.Fprintln(os.Stderr, "no file specified")
		os.Exit(2)
	}

	os.Exit(c.show(dir, filePath))
}

func (c *ShowCommand) show(dir *Directory, path string) int {
	realPath, err := dir.FilePath(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid path: %s\n", err)
		return 2
	}

	f, err := os.Open(realPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if fi.IsDir() {
		fmt.Fprintf(os.Stderr, "%s is a directory\n", path)
		return 1
	}

	_, err = io.Copy(os.Stdout, f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
