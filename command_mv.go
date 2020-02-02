package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type MvCommand struct {
	flag *flag.FlagSet

	force bool
}

func (c *MvCommand) Name() string {
	return "mv"
}

func (c *MvCommand) setup(args []string, options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		fmt.Fprintf(c.flag.Output(), "Usage: %s [options] <source> <destination>\n", c.Name())
		fmt.Fprintln(c.flag.Output(), "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.BoolVar(&c.force, "f", false, "allow overwriting existing files")
	c.flag.Parse(args)
}

func (c *MvCommand) Run(args []string, options *Options) {
	c.setup(args, options)

	dir, err := OpenDirectory(options.noteDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var src, dest string
	if c.flag.NArg() == 2 {
		src = c.flag.Arg(0)
		dest = c.flag.Arg(1)
	} else if c.flag.NArg() > 2 {
		fmt.Fprintln(os.Stderr, "too many arguments")
		os.Exit(2)
	} else {
		fmt.Fprintln(os.Stderr, "too few arguments")
		os.Exit(2)
	}

	c.move(dir, src, dest, options)
}

func (c *MvCommand) move(dir *Directory, src, dest string, options *Options) {
	srcReal, err := dir.FilePath(src)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	destReal, err := dir.FilePath(dest)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fi, err := os.Stat(destReal)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		if fi.IsDir() {
			destReal = filepath.Join(destReal, filepath.Base(srcReal))
		} else if !c.force {
			fmt.Fprintf(os.Stderr, "%s already exists; use -f to overwrite\n", dest)
			os.Exit(1)
		}
	}

	err = os.Rename(srcReal, destReal)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
