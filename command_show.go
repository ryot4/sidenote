package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

type ShowCommand struct {
	flag *flag.FlagSet
}

func (c *ShowCommand) Name() string {
	return "show"
}

func (c *ShowCommand) Description() string {
	return "Open notes with pager"
}

func (c *ShowCommand) setup(args []string, _options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		output := c.flag.Output()
		fmt.Fprintf(output, "Usage: %s %s <name>\n", os.Args[0], c.Name())
		fmt.Fprintf(output, "\n%s.\n", c.Description())
	}
	c.flag.Parse(args)
}

func (c *ShowCommand) Run(args []string, options *Options) error {
	c.setup(args, options)

	pager, ok := os.LookupEnv("PAGER")
	if !ok {
		return errors.New("PAGER is not set")
	}

	dir, err := checkDirectory(options.noteDir)
	if err != nil {
		return err
	}

	if c.flag.NArg() == 0 {
		return NewSyntaxError("no file specified")
	} else if c.flag.NArg() > 1 {
		return NewSyntaxError("too many arguments")
	}

	return c.runPager(dir, pager, c.flag.Arg(0))
}

func (c *ShowCommand) runPager(dir *Directory, pager, path string) error {
	realPath, err := dir.JoinPath(path)
	if err != nil {
		return err
	}

	pagerCmd := exec.Command(pager, realPath)
	pagerCmd.Stdin = os.Stdin
	pagerCmd.Stdout = os.Stdout
	pagerCmd.Stderr = os.Stderr
	return pagerCmd.Run()
}
