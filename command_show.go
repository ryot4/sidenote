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
		fmt.Fprintf(output, "Usage: %s %s <name>...\n", os.Args[0], c.Name())
		fmt.Fprintf(output, "\n%s.\n", c.Description())
	}
	c.flag.Parse(args)
}

func (c *ShowCommand) Run(args []string, options *Options) error {
	c.setup(args, options)

	d, err := FindStorageDir(options.NotesDirName())
	if err != nil {
		return err
	}
	storage := NewStorage(d)

	if c.flag.NArg() == 0 {
		return ErrNoFileName
	}

	return c.runPager(storage, c.flag.Args())
}

func (c *ShowCommand) runPager(storage *Storage, names []string) error {
	pager, ok := os.LookupEnv("PAGER")
	if !ok {
		return errors.New("PAGER is not set")
	}

	var args []string
	for _, name := range names {
		path, err := storage.JoinPath(name)
		if err != nil {
			return err
		}
		args = append(args, path)
	}

	pagerCmd := exec.Command(pager, args...)
	pagerCmd.Stdin = os.Stdin
	pagerCmd.Stdout = os.Stdout
	pagerCmd.Stderr = os.Stderr
	return pagerCmd.Run()
}
