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
		return err
	}

	err = os.MkdirAll(filepath.Dir(realPath), os.ModePerm)
	if err != nil {
		return err
	}

	fi, err := os.Stat(realPath)
	if os.IsNotExist(err) {
		if filepath.Ext(realPath) == "" && c.fileExt != "" {
			if strings.HasPrefix(c.fileExt, ".") {
				realPath += c.fileExt
			} else {
				realPath += "." + c.fileExt
			}
		}
	} else if err != nil {
		return err
	} else if fi.IsDir() {
		return fmt.Errorf("directory exists: %s", realPath)
	}

	editorCmd := exec.Command(editor, realPath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr
	return editorCmd.Run()
}
