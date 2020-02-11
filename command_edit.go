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
		fmt.Fprintf(c.flag.Output(), "Usage: %s [-f format] [name]\n", c.Name())
		fmt.Fprintln(c.flag.Output(), "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.StringVar(&c.nameFormat, "f", os.Getenv(NameFormatEnv),
		fmt.Sprintf("Generate filename using the given strftime format string (env: %s)",
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

	dir, err := checkDirectory(options.noteDir)
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
			exitWithSyntaxError("no filename specified")
		}
		filePath = Strftime(time.Now(), c.nameFormat)
	}

	err = c.runEditor(dir, editor, filePath)
	if err != nil {
		exitWithError(err)
	}
}

func (c *EditCommand) runEditor(dir *Directory, editor, path string) error {
	realPath, err := dir.JoinPath(path)
	if err != nil {
		exitWithError(err)
	}

	err = os.MkdirAll(filepath.Dir(realPath), os.ModePerm)
	if err != nil {
		return err
	}

	fi, err := os.Stat(realPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if fi.IsDir() {
		return fmt.Errorf("directory exists: %s", realPath)
	}

	editorCmd := exec.Command(editor, realPath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr
	return editorCmd.Run()
}
