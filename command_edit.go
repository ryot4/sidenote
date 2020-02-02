package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type EditCommand struct {
	flag *flag.FlagSet

	editor string
}

func (c *EditCommand) Name() string {
	return "edit"
}

func (c *EditCommand) setup(args []string, _options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		fmt.Fprintf(c.flag.Output(), "Usage: %s [options] path-to-file\n", c.Name())
		fmt.Fprintln(c.flag.Output(), "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.StringVar(&c.editor, "e", os.Getenv("EDITOR"), "editor to use")
	c.flag.Parse(args)
}

func (c *EditCommand) Run(args []string, options *Options) {
	c.setup(args, options)

	dir, err := OpenDirectory(options.noteDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var filePath string
	if c.flag.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "too many arguments")
		os.Exit(2)
	} else if c.flag.NArg() == 1 {
		filePath = c.flag.Arg(0)
	} else {
		fmt.Fprintln(os.Stderr, "no file specified")
		os.Exit(2)
	}
	realPath, err := dir.FilePath(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid path: %s\n", err)
		os.Exit(2)
	}
	runEditor(c.editor, realPath)
}

func runEditor(editor, path string) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "cannot create %s: %v\n", dir, err)
		os.Exit(1)
	}

	fi, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "cannot stat %s: %v\n", path, err)
			os.Exit(1)
		}
	} else if fi.IsDir() {
		fmt.Fprintf(os.Stderr, "directory exists: %s\n", path)
		os.Exit(1)
	}

	editorCmd := exec.Command(editor, path)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr
	editorCmd.Run()
	os.Exit(editorCmd.ProcessState.ExitCode())
}
