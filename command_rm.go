package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type RmCommand struct {
	flag *flag.FlagSet

	recurse bool
}

func (c *RmCommand) Name() string {
	return "rm"
}

func (c *RmCommand) setup(args []string, options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		fmt.Fprintf(c.flag.Output(), "Usage: %s [options] <name>\n", c.Name())
		fmt.Fprintln(c.flag.Output(), "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.BoolVar(&c.recurse, "r", false, "remove directories recursively")
	c.flag.Parse(args)
}

func (c *RmCommand) Run(args []string, options *Options) {
	c.setup(args, options)

	dir, err := OpenDirectory(options.noteDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var rmPath string
	if c.flag.NArg() == 1 {
		rmPath = c.flag.Arg(0)
	} else if c.flag.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "too many arguments")
		os.Exit(2)
	} else {
		fmt.Fprintln(os.Stderr, "too few arguments")
		os.Exit(2)
	}

	c.remove(dir, rmPath)
}

func (c *RmCommand) remove(dir *Directory, path string) {
	realPath, err := dir.FilePath(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fi, err := os.Stat(realPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if fi.IsDir() {
		isEmpty, err := isEmptyDir(realPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if !isEmpty && !c.recurse {
			fmt.Fprintln(os.Stderr, "directory not empty: use -r to remove")
			os.Exit(1)
		}
	}
	err = os.RemoveAll(realPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func isEmptyDir(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
