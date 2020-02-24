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

func (c *PathCommand) setup(args []string, _options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		fmt.Fprintf(c.flag.Output(), "Usage: %s [-a] [name]\n", c.Name())
		fmt.Fprintln(c.flag.Output(), "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.BoolVar(&c.followSymlink, "L", false, "Follow the symbolic link to notes (implies -a when the target path is absolute)")
	c.flag.BoolVar(&c.absolute, "a", false, "Show absolute path")
	c.flag.BoolVar(&c.check, "c", false, "Check existence of the path")
	c.flag.Parse(args)
}

func (c *PathCommand) Run(args []string, options *Options) {
	c.setup(args, options)

	dir, err := getDirectory(options.noteDir)
	if err != nil {
		exitWithError(err)
	}

	if c.flag.NArg() > 1 {
		exitWithSyntaxError("too many arguments")
	}
	path := c.flag.Arg(0)
	if path == "" {
		path = "/"
	}

	err = c.showPath(dir, path)
	if err != nil {
		exitWithError(err)
	}
}

func (c *PathCommand) showPath(dir *Directory, path string) error {
	if c.followSymlink {
		link, err := dir.FollowSymlink()
		if err != nil {
			return err
		}
		if link && dir.IsAbs() {
			c.absolute = true
		}
	}

	realPath, err := dir.JoinPath(path)
	if err != nil {
		return err
	}

	if c.check {
		_, err := os.Stat(realPath)
		if err != nil {
			return err
		}
	}

	if c.absolute {
		absPath, err := filepath.Abs(realPath)
		if err != nil {
			return err
		}
		fmt.Println(absPath)
	} else {
		if filepath.IsAbs(realPath) {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			relPath, err := filepath.Rel(wd, realPath)
			if err != nil {
				return err
			}
			fmt.Println(relPath)
		} else {
			fmt.Println(realPath)
		}
	}
	return nil
}
