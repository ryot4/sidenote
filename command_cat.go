package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type CatCommand struct {
	flag *flag.FlagSet
}

func (c *CatCommand) Name() string {
	return "cat"
}

func (c *CatCommand) Description() string {
	return "Print contents of notes"
}

func (c *CatCommand) setup(args []string, _options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		output := c.flag.Output()
		fmt.Fprintf(output, "Usage: %s <name>...\n", c.Name())
		fmt.Fprintf(output, "\n%s.\n", c.Description())
	}
	c.flag.Parse(args)
}

func (c *CatCommand) Run(args []string, options *Options) {
	c.setup(args, options)

	dir, err := checkDirectory(options.noteDir)
	if err != nil {
		exitWithError(err)
	}

	if c.flag.NArg() == 0 {
		exitWithSyntaxError("no file specified")
	}

	var lastErr error
	for _, filePath := range c.flag.Args() {
		err = c.catFile(dir, filePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			lastErr = err
		}
	}
	if lastErr != nil {
		os.Exit(1)
	}
}

func (c *CatCommand) catFile(dir *Directory, path string) error {
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
