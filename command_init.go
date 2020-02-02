package main

import (
	"flag"
	"fmt"
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
		_, err = InitDirectory(options.noteDir)
	} else {
		_, err = InitDirectoryLink(c.linkDir, options.noteDir)
	}
	if err != nil {
		exitWithError(err)
	}
}
