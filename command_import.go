package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ImportCommand struct {
	flag *flag.FlagSet

	force  bool
	delete bool
}

func (c *ImportCommand) Name() string {
	return "import"
}

func (c *ImportCommand) Description() string {
	return "Import a note from the existing file or the standard input"
}

func (c *ImportCommand) setup(args []string, options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		output := c.flag.Output()
		fmt.Fprintf(output, "Usage: %s %s [-d] [-f] {<file>|-} [<name>]\n", os.Args[0], c.Name())
		fmt.Fprintf(output, "\n%s.\n", c.Description())
		fmt.Fprintln(output, "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.BoolVar(&c.delete, "d", false, "Delete the original file after import")
	c.flag.BoolVar(&c.force, "f", false, "Allow overwriting existing files")
	c.flag.Parse(args)
}

func (c *ImportCommand) Run(args []string, options *Options) {
	c.setup(args, options)

	dir, err := checkDirectory(options.noteDir)
	if err != nil {
		exitWithError(err)
	}

	var origPath, importedPath string
	switch c.flag.NArg() {
	case 0:
		exitWithSyntaxError("no file specified")
	case 1:
		origPath = c.flag.Arg(0)
		if origPath == "-" {
			exitWithSyntaxError("no name specified")
		}
		importedPath = filepath.Base(origPath)
	case 2:
		origPath = c.flag.Arg(0)
		importedPath = c.flag.Arg(1)
	default:
		exitWithSyntaxError("too many arguments")
	}

	var r io.Reader
	if origPath == "-" {
		r = os.Stdin
	} else {
		file, err := os.Open(origPath)
		if err != nil {
			exitWithError(err)
		}
		defer file.Close()

		fi, err := file.Stat()
		if err != nil {
			exitWithError(err)
		} else if !fi.Mode().IsRegular() {
			exitWithError(fmt.Errorf("%s is not a regular file", origPath))
		}
		r = file
	}

	err = c.importFile(dir, importedPath, r)
	if err != nil {
		exitWithError(err)
	}

	if c.delete {
		if origPath == "-" {
			fmt.Fprintln(os.Stderr, "ignoring -d; the input is from the standard input")
		} else {
			err = os.Remove(origPath)
			if err != nil {
				exitWithError(err)
			}
		}
	}
}

func (c *ImportCommand) importFile(dir *Directory, path string, r io.Reader) error {
	realPath, err := dir.JoinPath(path)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(realPath), os.ModePerm)
	if err != nil {
		return err
	}

	fi, err := os.Stat(realPath)
	if err == nil {
		if fi.IsDir() {
			return fmt.Errorf("directory exists: %s", realPath)
		} else if !c.force {
			return fmt.Errorf("file exists; use -f to overwrite: %s", realPath)
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	file, err := os.Create(realPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return err
}
