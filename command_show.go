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
		fmt.Fprintf(c.flag.Output(), "Usage: %s <name>\n", c.Name())
	}
	c.flag.Parse(args)
}

func (c *ShowCommand) Run(args []string, options *Options) {
	c.setup(args, options)

	dir, err := checkDirectory(options.noteDir)
	if err != nil {
		exitWithError(err)
	}

	var filePath string
	if c.flag.NArg() > 1 {
		exitWithSyntaxError("too many arguments")
	} else if c.flag.NArg() == 1 {
		filePath = c.flag.Arg(0)
	} else {
		exitWithSyntaxError("no file specified")
	}

	err = c.show(dir, filePath)
	if err != nil {
		exitWithError(err)
	}
}

func (c *ShowCommand) show(dir *Directory, path string) error {
	realPath, err := dir.JoinPath(path)
	if err != nil {
		return err
	}

	f, err := os.Open(realPath)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return fmt.Errorf("%s is a directory", path)
	}

	_, err = io.Copy(os.Stdout, f)
	if err != nil {
		return err
	}
	return nil
}
