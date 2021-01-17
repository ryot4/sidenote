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
		fmt.Fprintf(output, "Usage: %s %s [-r] <name>...\n", os.Args[0], c.Name())
		fmt.Fprintf(output, "\n%s.\n", c.Description())
		fmt.Fprintln(output, "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.BoolVar(&c.recurse, "r", false, "Remove directories recursively")
	c.flag.Parse(args)
}

func (c *RmCommand) Run(args []string, options *Options) error {
	c.setup(args, options)

	dir, err := options.CheckDirectory()
	if err != nil {
		return err
	}

	if c.flag.NArg() == 0 {
		return ErrNoFileName
	}

	for _, name := range c.flag.Args() {
		err = c.remove(dir, name)
		if err != nil {
			return err
		}
		fmt.Printf("removed %s\n", name)
	}
	return nil
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
			return fmt.Errorf("directory %s is not empty (use -r to remove recursively)", name)
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
