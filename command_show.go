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

func (c *ShowCommand) setup(args []string, _options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		fmt.Fprintf(c.flag.Output(), "Usage: %s <name>\n", c.Name())
	}
	c.flag.Parse(args)
}

func (c *ShowCommand) Run(args []string, options *Options) {
	c.setup(args, options)

	pager, ok := os.LookupEnv("PAGER")
	if !ok {
		exitWithError(errors.New("PAGER is not set"))
	}

	dir, err := checkDirectory(options.noteDir)
	if err != nil {
		exitWithError(err)
	}

	if c.flag.NArg() == 0 {
		exitWithSyntaxError("no file specified")
	} else if c.flag.NArg() > 1 {
		exitWithSyntaxError("too many arguments")
	}

	err = c.runPager(dir, pager, c.flag.Arg(0))
	if err != nil {
		exitWithError(err)
	}
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
