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
	return "Initialize the directory for notes"
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

func (c *InitCommand) Run(args []string, options *Options) error {
	c.setup(args, options)

	if c.flag.NArg() > 0 {
		return NewSyntaxError("too many arguments")
	}

	noteDir := options.NoteDirName()

	if c.linkTarget == "" {
		err := NewDirectory(noteDir).Init()
		if err != nil {
			return err
		}
		fmt.Printf("initialized %s\n", noteDir)
	} else {
		// Create the symlink first; if this fails, do not initialize the target directory.
		err := os.Symlink(c.linkTarget, noteDir)
		if err != nil {
			return err
		}

		err = NewDirectory(c.linkTarget).Init()
		if err != nil && !os.IsExist(err) {
			return err
		}
		fmt.Printf("initialized %s (-> %s)\n", noteDir, c.linkTarget)
	}
	return nil
}
