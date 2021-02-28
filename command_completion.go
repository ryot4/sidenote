package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
)

//go:embed completions/sidenote.bash
var bashCompletionScript string

type CompletionCommand struct {
	flag *flag.FlagSet
}

func (c *CompletionCommand) Name() string {
	return "completion"
}

func (c *CompletionCommand) Description() string {
	return "Print shell function for command line completion"
}

func (c *CompletionCommand) setup(args []string, options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		output := c.flag.Output()
		fmt.Fprintf(output, "Usage: %s %s <shell>\n", os.Args[0], c.Name())
		fmt.Fprintf(output, "\n%s.\n", c.Description())
	}
	c.flag.Parse(args)
}

func (c *CompletionCommand) Run(args []string, options *Options) error {
	c.setup(args, options)

	if c.flag.NArg() == 0 {
		return NewSyntaxError("no shell name specified")
	}

	shellName := c.flag.Arg(0)
	switch shellName {
	case "bash":
		fmt.Print(bashCompletionScript)
	default:
		return NewSyntaxError(fmt.Sprintf("Unsupported shell: %s", shellName))
	}
	return nil
}
