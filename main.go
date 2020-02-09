package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	NoteDirName = ".notes"
	NoteDirEnv  = "SIDENOTE_DIR"
)

type Options struct {
	noteDir string
}

type Command interface {
	Name() string
	Run(args []string, options *Options)
}

var subCommands = []Command{
	&InitCommand{},
	&EditCommand{},
	&LsCommand{},
	&ShowCommand{},
	&MvCommand{},
	&RmCommand{},
}

func main() {
	var options Options
	var printVersion bool

	flag.Usage = usage
	flag.StringVar(&options.noteDir, "d", os.Getenv(NoteDirEnv),
		fmt.Sprintf("Specify the directory for notes (env: %s)", NoteDirEnv))
	flag.BoolVar(&printVersion, "version", false, "Print the version and exit")
	flag.Parse()

	if printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	run(flag.Args(), &options)
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <command> [command-arguments]\n", os.Args[0])
	fmt.Fprintln(flag.CommandLine.Output(), "\noptions:")
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\ncommands (run \"%s <command> -h\" for command usage):\n", os.Args[0])
	for _, cmd := range subCommands {
		fmt.Fprintf(flag.CommandLine.Output(), "  %s\n", cmd.Name())
	}
}

func run(args []string, options *Options) {
	if len(args) > 0 {
		cmdName := args[0]
		for _, cmd := range subCommands {
			if cmdName == cmd.Name() {
				cmd.Run(args[1:], options)
				return
			}
		}
		exitWithSyntaxError(fmt.Sprintf("unknown command %q", cmdName))
	} else {
		exitWithSyntaxError("no command specified")
	}
}
