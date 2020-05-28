package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func (c *ImportCommand) Run(args []string, options *Options) error {
	c.setup(args, options)

	dir, err := checkDirectory(options.noteDir)
	if err != nil {
		return err
	}

	var origPath, name string
	switch c.flag.NArg() {
	case 0:
		return NewSyntaxError("no file specified")
	case 1:
		origPath = c.flag.Arg(0)
		if origPath == "-" {
			return NewSyntaxError("no name specified (required when importing from the standard input)")
		}
		name = filepath.Base(origPath)
	case 2:
		origPath = c.flag.Arg(0)
		name = c.flag.Arg(1)
		if strings.HasSuffix(name, string(filepath.Separator)) {
			if origPath == "-" {
				return NewSyntaxError("no name specified (required when importing from the standard input)")
			} else {
				name = filepath.Join(name, filepath.Base(origPath))
			}
		}
	default:
		return NewSyntaxError("too many arguments")
	}

	var r io.Reader
	if origPath == "-" {
		r = os.Stdin
	} else {
		file, err := os.Open(origPath)
		if err != nil {
			return err
		}
		defer file.Close()

		fi, err := file.Stat()
		if err != nil {
			return err
		} else if !fi.Mode().IsRegular() {
			return fmt.Errorf("%s is not a regular file", origPath)
		}
		r = file
	}

	err = c.importFile(dir, name, r)
	if err != nil {
		return err
	}

	if c.delete {
		if origPath == "-" {
			fmt.Fprintln(os.Stderr, "ignoring -d; the input is from the standard input")
		} else {
			err = os.Remove(origPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *ImportCommand) importFile(dir *Directory, name string, r io.Reader) error {
	path, err := dir.JoinPath(name)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return fmt.Errorf("directory exists (use %s/ to import into the directory): %s", name, name)
		} else if !c.force {
			return fmt.Errorf("file exists; use -f to overwrite: %s", name)
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return err
}
