package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type ExecCommand struct {
	flag *flag.FlagSet

	chdir string
}

func (c *ExecCommand) Name() string {
	return "exec"
}

func (c *ExecCommand) Description() string {
	return "Execute a command inside notes directory"
}

func (c *ExecCommand) setup(args []string, _options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		output := c.flag.Output()
		fmt.Fprintf(output, "Usage: %s %s [-cd dir] command\n", os.Args[0], c.Name())
		fmt.Fprintf(output, "\n%s.\n", c.Description())
		fmt.Fprintln(output, "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.StringVar(&c.chdir, "cd", "", "Change the directory where commands are executed")
	c.flag.Parse(args)
}

func (c *ExecCommand) Run(args []string, options *Options) error {
	c.setup(args, options)

	d, err := FindStorageDir(options.NotesDirName())
	if err != nil {
		return err
	}
	storage := NewStorage(d)

	if c.flag.NArg() == 0 {
		return NewSyntaxError("no command specified")
	}
	return c.execCommand(storage, c.flag.Args())
}

func (c *ExecCommand) execCommand(storage *Storage, command []string) error {
	path, err := storage.JoinPath(c.chdir)
	if err != nil {
		return err
	}

	err = os.Chdir(path)
	if err != nil {
		return err
	}

	program, err := exec.LookPath(command[0])
	if err != nil {
		return err
	}
	return syscall.Exec(program, command, os.Environ())
}
