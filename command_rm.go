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

func (c *RmCommand) Description() string {
	return "Remove notes"
}

func (c *RmCommand) setup(args []string, options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		fmt.Fprintf(c.flag.Output(), "Usage: %s [-r] <name>\n", c.Name())
		fmt.Fprintf(c.flag.Output(), "\n%s.\n", c.Description())
		fmt.Fprintln(c.flag.Output(), "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.BoolVar(&c.recurse, "r", false, "Remove directories recursively")
	c.flag.Parse(args)
}

func (c *RmCommand) Run(args []string, options *Options) {
	c.setup(args, options)

	dir, err := checkDirectory(options.noteDir)
	if err != nil {
		exitWithError(err)
	}

	var rmPath string
	if c.flag.NArg() == 1 {
		rmPath = c.flag.Arg(0)
	} else if c.flag.NArg() > 1 {
		exitWithSyntaxError("too many arguments")
	} else {
		exitWithSyntaxError("too few arguments")
	}

	err = c.remove(dir, rmPath)
	if err != nil {
		exitWithError(err)
	}
}

func (c *RmCommand) remove(dir *Directory, path string) error {
	realPath, err := dir.JoinPath(path)
	if err != nil {
		return err
	}
	fi, err := os.Stat(realPath)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		isEmpty, err := isEmptyDir(realPath)
		if err != nil {
			return err
		}
		if !isEmpty && !c.recurse {
			return fmt.Errorf("directory not empty: use -r to remove")
		}
	}
	return os.RemoveAll(realPath)
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
