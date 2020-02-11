package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type MvCommand struct {
	flag *flag.FlagSet

	force bool
}

func (c *MvCommand) Name() string {
	return "mv"
}

func (c *MvCommand) setup(args []string, options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		fmt.Fprintf(c.flag.Output(), "Usage: %s [-f] <source> <destination>\n", c.Name())
		fmt.Fprintln(c.flag.Output(), "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.BoolVar(&c.force, "f", false, "Allow overwriting existing files")
	c.flag.Parse(args)
}

func (c *MvCommand) Run(args []string, options *Options) {
	c.setup(args, options)

	dir, err := checkDirectory(options.noteDir)
	if err != nil {
		exitWithError(err)
	}

	var src, dest string
	if c.flag.NArg() == 2 {
		src = c.flag.Arg(0)
		dest = c.flag.Arg(1)
	} else if c.flag.NArg() > 2 {
		exitWithSyntaxError("too many arguments")
	} else {
		exitWithSyntaxError("too few arguments")
	}

	err = c.move(dir, src, dest)
	if err != nil {
		exitWithError(err)
	}
}

func (c *MvCommand) move(dir *Directory, src, dest string) error {
	srcReal, err := dir.JoinPath(src)
	if err != nil {
		return err
	}
	destReal, err := dir.JoinPath(dest)
	if err != nil {
		return err
	}

	fi, err := os.Stat(destReal)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if fi.IsDir() {
		destReal = filepath.Join(destReal, filepath.Base(srcReal))
	}

	return c.doMove(srcReal, destReal)
}

func (c *MvCommand) doMove(srcReal, destReal string) error {
	_, err := os.Stat(destReal)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if !c.force {
		return fmt.Errorf("%s already exists; use -f to overwrite", destReal)
	}

	err = os.MkdirAll(filepath.Dir(destReal), os.ModePerm)
	if err != nil {
		return err
	}
	return os.Rename(srcReal, destReal)
}
