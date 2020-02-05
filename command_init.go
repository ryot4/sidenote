package main

import (
	"flag"
	"fmt"
	"os"
)

type InitCommand struct {
	flag *flag.FlagSet

	linkDir string
}

func (c *InitCommand) Name() string {
	return "init"
}

func (c *InitCommand) setup(args []string, options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		fmt.Fprintf(c.flag.Output(), "Usage: %s [options]\n", c.Name())
		fmt.Fprintln(c.flag.Output(), "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.StringVar(&c.linkDir, "l", "", fmt.Sprintf("make %s a symlink to the directory", options.noteDir))
	c.flag.Parse(args)
}

func (c *InitCommand) Run(args []string, options *Options) {
	c.setup(args, options)

	if c.flag.NArg() > 0 {
		exitWithSyntaxError("too many arguments")
	}

	var err error
	if c.linkDir == "" {
		err = c.initDirectory(options.noteDir)
	} else {
		err = c.initLink(options.noteDir)
	}
	if err != nil {
		exitWithError(err)
	}
}

func (c *InitCommand) initDirectory(noteDir string) error {
	_, err := InitDirectory(noteDir)
	return err
}

func (c *InitCommand) initLink(noteDir string) error {
	_, err := InitDirectory(c.linkDir)
	if err != nil {
		return err
	}
	return os.Symlink(c.linkDir, noteDir)
}
