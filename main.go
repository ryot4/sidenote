package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	NoteDirName = ".notes"
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
	&MvCommand{},
	&RmCommand{},
}

func main() {
	var options Options

	flag.Usage = usage
	flag.StringVar(&options.noteDir, "d", findNoteDir(), "path to the directory for notes")
	flag.Parse()

	run(flag.Args(), &options)
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <command> [command-args]\n", os.Args[0])
	fmt.Fprintln(flag.CommandLine.Output(), "\noptions:")
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\ncommands (run \"%s <command> -h\" for command usage):\n", os.Args[0])
	for _, cmd := range subCommands {
		fmt.Fprintf(flag.CommandLine.Output(), "  %s\n", cmd.Name())
	}
}

func findNoteDir() string {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot get current directory: %s\n", err)
		os.Exit(1)
	}
	separator := string(filepath.Separator)
	for dir := wd; dir != "." && dir != separator; dir = filepath.Dir(dir) {
		noteDir := filepath.Join(dir, NoteDirName)
		fi, err := os.Stat(noteDir)
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "cannot stat %s: %s. ignoring\n", noteDir, err)
			}
			continue
		}
		if fi.IsDir() {
			rel, err := filepath.Rel(wd, noteDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot get relative path to %s: %s. ignoring\n", noteDir, err)
				continue
			}
			return rel
		}
	}
	return NoteDirName
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
		fmt.Fprintf(os.Stderr, "unknown command %q\n", cmdName)
		os.Exit(2)
	} else {
		fmt.Fprintln(os.Stderr, "no command specified")
		os.Exit(2)
	}
}
