package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	NameFormatEnv = "SIDENOTE_NAME_FORMAT"
	FileExtEnv    = "SIDENOTE_FILE_EXT"
)

type EditCommand struct {
	flag *flag.FlagSet

	nameFormat string
	fileExt    string
}

func (c *EditCommand) Name() string {
	return "edit"
}

func (c *EditCommand) Description() string {
	return "Open notes with text editor"
}

func (c *EditCommand) setup(args []string, _options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		output := c.flag.Output()
		fmt.Fprintf(output, "Usage: %s %s [-f format] [-x extension] [name]\n", os.Args[0], c.Name())
		fmt.Fprintf(output, "\n%s.\n", c.Description())
		fmt.Fprintln(output, "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.StringVar(&c.nameFormat, "f", os.Getenv(NameFormatEnv),
		fmt.Sprintf("Generate filename using the given strftime format string (env: %s)",
			NameFormatEnv))
	c.flag.StringVar(&c.fileExt, "x", os.Getenv(FileExtEnv),
		fmt.Sprintf("Specify the default file extension for new files (env: %s)",
			FileExtEnv))
	c.flag.Parse(args)
}

func (c *EditCommand) Run(args []string, options *Options) error {
	c.setup(args, options)

	dir, err := checkDirectory(options.noteDir)
	if err != nil {
		return err
	}

	var name string
	switch c.flag.NArg() {
	case 0:
		if c.nameFormat == "" {
			return NewSyntaxError("no filename specified")
		}
		name = Strftime(time.Now(), c.nameFormat)
	case 1:
		name = c.flag.Arg(0)
	default:
		return NewSyntaxError("too many arguments")
	}

	return c.runEditor(dir, name)
}

func (c *EditCommand) runEditor(dir *Directory, name string) error {
	path, err := dir.JoinPath(name)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		if filepath.Ext(path) == "" && c.fileExt != "" {
			if strings.HasPrefix(c.fileExt, ".") {
				path += c.fileExt
			} else {
				path += "." + c.fileExt
			}
		}
	} else if err != nil {
		return err
	} else if fi.IsDir() {
		return fmt.Errorf("directory exists: %s", name)
	}

	editor, ok := os.LookupEnv("VISUAL")
	if !ok {
		editor, ok = os.LookupEnv("EDITOR")
		if !ok {
			return errors.New("neither VISUAL nor EDITOR is set")
		}
	}

	editorCmd := exec.Command(editor, path)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr
	return editorCmd.Run()
}
