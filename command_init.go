package main

import (
	"flag"
	"fmt"
	"os"
)

type InitCommand struct {
	flag *flag.FlagSet

	linkTarget string
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
	c.flag.StringVar(&c.linkTarget, "l", "", fmt.Sprintf("Link %s to the specified directory", NoteDirName))
	c.flag.Parse(args)
}

func (c *InitCommand) Run(args []string, options *Options) {
	c.setup(args, options)

	if c.flag.NArg() > 0 {
		exitWithSyntaxError("too many arguments")
	}

	var err error
	if c.linkTarget == "" {
		err = c.initDirectory(options)
	} else {
		err = c.initLink(options)
	}
	if err != nil {
		exitWithError(err)
	}
}

func (c *InitCommand) initDirectory(options *Options) error {
	noteDir := NoteDirName
	if options.noteDir != "" {
		noteDir = options.noteDir
	}
	_, err := InitDirectory(noteDir)
	return err
}

func (c *InitCommand) initLink(options *Options) error {
	if options.noteDir != "" {
		fmt.Fprintln(os.Stderr, "warning: -d is ignored when -l is specified")
	}
	_, err := InitDirectory(c.linkTarget)
	if err != nil {
		return err
	}
	return os.Symlink(c.linkTarget, NoteDirName)
}
