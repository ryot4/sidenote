package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type PathCommand struct {
	flag *flag.FlagSet

	absolute bool
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
	c.flag.BoolVar(&c.absolute, "a", false, "Show absolute path")
	c.flag.Parse(args)
}

func (c *PathCommand) Run(args []string, options *Options) {
	c.setup(args, options)

	dir, err := openDirectory(options.noteDir)
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
	absPath, err := dir.FilePath(path)
	if err != nil {
		exitWithError(err)
	}

	if c.absolute {
		fmt.Println(absPath)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(wd, absPath)
		if err != nil {
			return err
		}
		fmt.Println(relPath)
	}
	return nil
}
