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
	return "Open notes with $PAGER"
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

	dir, err := options.CheckDirectory()
	if err != nil {
		return err
	}

	var name string
	switch c.flag.NArg() {
	case 0:
		return ErrNoFileName
	case 1:
		name = c.flag.Arg(0)
	default:
		return ErrTooManyArgs
	}

	return c.runPager(dir, name)
}

func (c *ShowCommand) runPager(dir *Directory, name string) error {
	path, err := dir.JoinPath(name)
	if err != nil {
		return err
	}

	pager, ok := os.LookupEnv("PAGER")
	if !ok {
		return errors.New("PAGER is not set")
	}

	pagerCmd := exec.Command(pager, path)
	pagerCmd.Stdin = os.Stdin
	pagerCmd.Stdout = os.Stdout
	pagerCmd.Stderr = os.Stderr
	return pagerCmd.Run()
}
