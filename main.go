package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
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
	flag.StringVar(&options.noteDir, "d", findNoteDir(), "Specify the directory for notes")
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

func findNoteDir() string {
	envDir := os.Getenv(NoteDirEnv)
	if envDir != "" {
		return envDir
	}

	// Search ".notes" directory upward from the current directory
	wd, err := os.Getwd()
	if err != nil {
		exitWithError(err)
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
		exitWithSyntaxError(fmt.Sprintf("unknown command %q", cmdName))
	} else {
		exitWithSyntaxError("no command specified")
	}
}
