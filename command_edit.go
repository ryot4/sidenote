package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	NameFormatEnv = "SIDENOTE_NAME_FORMAT"
)

type EditCommand struct {
	flag *flag.FlagSet

	nameFormat string
}

func (c *EditCommand) Name() string {
	return "edit"
}

func (c *EditCommand) setup(args []string, _options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		fmt.Fprintf(c.flag.Output(), "Usage: %s [options] [name]\n", c.Name())
		fmt.Fprintln(c.flag.Output(), "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.StringVar(&c.nameFormat, "f", os.Getenv(NameFormatEnv),
		fmt.Sprintf("Generate file name using the given strftime format string (env: %s)",
			NameFormatEnv))
	c.flag.Parse(args)
}

func (c *EditCommand) Run(args []string, options *Options) {
	c.setup(args, options)

	editor, ok := os.LookupEnv("VISUAL")
	if !ok {
		editor, ok = os.LookupEnv("EDITOR")
		if !ok {
			exitWithError(errors.New("neither VISUAL nor EDITOR is set"))
		}
	}

	dir, err := openDirectory(options.noteDir)
	if err != nil {
		exitWithError(err)
	}

	var filePath string
	if c.flag.NArg() > 1 {
		exitWithSyntaxError("too many arguments")
	} else if c.flag.NArg() == 1 {
		filePath = c.flag.Arg(0)
	} else {
		if c.nameFormat == "" {
			exitWithSyntaxError("no file name specified")
		}
		filePath = Strftime(time.Now(), c.nameFormat)
	}
	absPath, err := dir.AbsPath(filePath)
	if err != nil {
		exitWithError(err)
	}
	err = runEditor(editor, absPath)
	if err != nil {
		exitWithError(err)
	}
}

func runEditor(editor, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	fi, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if fi.IsDir() {
		return fmt.Errorf("directory exists: %s", path)
	}

	editorCmd := exec.Command(editor, path)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr
	return editorCmd.Run()
}
