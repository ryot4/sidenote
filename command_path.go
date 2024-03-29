package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type PathCommand struct {
	flag *flag.FlagSet

	absolute      bool
	check         bool
	followSymlink bool
}

func (c *PathCommand) Name() string {
	return "path"
}

func (c *PathCommand) Description() string {
	return "Print the path of notes"
}

func (c *PathCommand) setup(args []string, _options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		output := c.flag.Output()
		fmt.Fprintf(output, "Usage: %s %s [-L] [-a] [-c] [name]\n", os.Args[0], c.Name())
		fmt.Fprintf(output, "\n%s.\n", c.Description())
		fmt.Fprintln(output, "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.BoolVar(&c.followSymlink, "L", false, "Follow the symbolic link to notes (implies -a when the target path is absolute)")
	c.flag.BoolVar(&c.absolute, "a", false, "Show absolute path")
	c.flag.BoolVar(&c.check, "c", false, "Check existence of the path")
	c.flag.Parse(args)
}

func (c *PathCommand) Run(args []string, options *Options) error {
	c.setup(args, options)

	d, err := options.FindDirectory()
	if err != nil {
		return err
	}

	if c.followSymlink {
		isLink, err := isSymlink(d)
		if err != nil {
			return err
		}
		if isLink {
			d, err = os.Readlink(d)
			if err != nil {
				return err
			}
			if filepath.IsAbs(d) {
				c.absolute = true
			}
		}
	}

	if c.flag.NArg() > 1 {
		return ErrTooManyArgs
	}
	name := c.flag.Arg(0)
	if name == "" {
		name = "/"
	}

	return c.showPath(NewDirectory(d), name)
}

func isSymlink(path string) (bool, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return fi.Mode()&os.ModeSymlink != 0, nil
}

func (c *PathCommand) showPath(dir *Directory, name string) error {
	path, err := dir.JoinPath(name)
	if err != nil {
		return err
	}

	if c.check {
		_, err := os.Stat(path)
		if err != nil {
			return err
		}
	}

	if c.absolute {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		fmt.Println(absPath)
	} else {
		if filepath.IsAbs(path) {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			relPath, err := filepath.Rel(wd, path)
			if err != nil {
				return err
			}
			fmt.Println(relPath)
		} else {
			fmt.Println(path)
		}
	}
	return nil
}
