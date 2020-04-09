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

func (c *InitCommand) Description() string {
	return "Initialize a directory for notes"
}

func (c *InitCommand) setup(args []string, options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		output := c.flag.Output()
		fmt.Fprintf(output, "Usage: %s %s [-l path]\n", os.Args[0], c.Name())
		fmt.Fprintf(output, "\n%s.\n", c.Description())
		fmt.Fprintln(output, "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.StringVar(&c.linkTarget, "l", "", "Link notes to the specified directory")
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
	return NewDirectory(noteDir).Init()
}

func (c *InitCommand) initLink(options *Options) error {
	noteDir := NoteDirName
	if options.noteDir != "" {
		noteDir = options.noteDir
	}

	// Create the symlink first; if this fails, do not initialize the target directory.
	err := os.Symlink(c.linkTarget, noteDir)
	if err != nil {
		return err
	}

	err = NewDirectory(c.linkTarget).Init()
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}
