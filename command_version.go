package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "0.1.9"

type VersionCommand struct {
	flag *flag.FlagSet
}

func (c *VersionCommand) Name() string {
	return "version"
}

func (c *VersionCommand) Description() string {
	return "Print the version"
}

func (c *VersionCommand) setup(args []string, options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		output := c.flag.Output()
		fmt.Fprintf(output, "Usage: %s %s\n", os.Args[0], c.Name())
		fmt.Fprintf(output, "\n%s.\n", c.Description())
	}
	c.flag.Parse(args)
}

func (c *VersionCommand) Run(args []string, options *Options) error {
	c.setup(args, options)
	fmt.Println(version)
	return nil
}
