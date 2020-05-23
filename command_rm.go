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
		output := c.flag.Output()
		fmt.Fprintf(output, "Usage: %s %s [-r] <name>\n", os.Args[0], c.Name())
		fmt.Fprintf(output, "\n%s.\n", c.Description())
		fmt.Fprintln(output, "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.BoolVar(&c.recurse, "r", false, "Remove directories recursively")
	c.flag.Parse(args)
}

func (c *RmCommand) Run(args []string, options *Options) error {
	c.setup(args, options)

	dir, err := checkDirectory(options.noteDir)
	if err != nil {
		return err
	}

	var name string
	switch c.flag.NArg() {
	case 0:
		return NewSyntaxError("too few arguments")
	case 1:
		name = c.flag.Arg(0)
	default:
		return NewSyntaxError("too many arguments")
	}

	return c.remove(dir, name)
}

func (c *RmCommand) remove(dir *Directory, name string) error {
	path, err := dir.JoinPath(name)
	if err != nil {
		return err
	}
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		isEmpty, err := isEmptyDir(path)
		if err != nil {
			return err
		}
		if !isEmpty && !c.recurse {
			return fmt.Errorf("directory not empty: use -r to remove")
		}
	}
	return os.RemoveAll(path)
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
