package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

const NotesDirEnv = "SIDENOTE_DIR"

type Options struct {
	notesDir string
}

func (options *Options) NotesDirName() string {
	if options.notesDir != "" {
		return options.notesDir
	}
	return ".notes"
}

func (options *Options) CheckDirectory() (dir *Directory, err error) {
	dir, err = options.FindDirectory()
	if err != nil {
		return
	}
	isDir, err := dir.IsDir()
	if err != nil {
		return nil, err
	} else if !isDir {
		return nil, fmt.Errorf("%s exists but is not a directory", options.notesDir)
	}
	return
}

func (options *Options) FindDirectory() (dir *Directory, err error) {
	return FindDirectory(options.NotesDirName())
}

func main() {
	var options Options
	var printVersion bool

	flag.Usage = usage
	flag.StringVar(&options.notesDir, "d", "",
		fmt.Sprintf("Specify the directory for notes (env: %s)", NotesDirEnv))
	flag.BoolVar(&printVersion, "version", false, "Print the version and exit")
	flag.Parse()

	if options.notesDir == "" {
		options.notesDir = os.Getenv(NotesDirEnv)
	}

	if printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	run(flag.Args(), &options)
}

func usage() {
	output := flag.CommandLine.Output()
	fmt.Fprintf(output, "Usage: %s [-d path] [-version] <command> [command-arguments]\n", os.Args[0])
	fmt.Fprintln(output, "\noptions:")
	flag.PrintDefaults()
	fmt.Fprintf(output, "\ncommands:\n")
	w := tabwriter.NewWriter(output, 0, 0, 4, ' ', 0)
	for _, cmd := range subCommands {
		fmt.Fprintf(w, "  %s\t%s\n", cmd.Name(), cmd.Description())
	}
	w.Flush()
	fmt.Fprintf(output, "\nRun %s <command> -h for usage of each command.\n", os.Args[0])
}

var subCommands = []Command{
	&CatCommand{},
	&EditCommand{},
	&ExecCommand{},
	&ImportCommand{},
	&InitCommand{},
	&LsCommand{},
	&PathCommand{},
	&RmCommand{},
	&ServeCommand{},
	&ShowCommand{},
}

func run(args []string, options *Options) {
	if len(args) > 0 {
		cmdName := args[0]
		for _, cmd := range subCommands {
			if cmdName == cmd.Name() {
				err := cmd.Run(args[1:], options)
				if err != nil {
					exitWithError(err)
				}
				return
			}
		}
		exitWithError(NewSyntaxError(fmt.Sprintf("unknown command %q", cmdName)))
	} else {
		exitWithError(NewSyntaxError("no command specified"))
	}
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	if _, ok := err.(*SyntaxError); ok {
		os.Exit(2)
	} else {
		os.Exit(1)
	}
}
