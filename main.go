package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	NoteDir = ".notes"
)

type HandlerFunc func(notePath string, args []string)

var handlers = map[string]HandlerFunc{
	"init": runInit,
	"edit": runEdit,
	"ls":   runLs,
}

func main() {
	notePath := flag.String("d", findNotePath(), "path to the directory for notes")
	flag.Parse()
	run(*notePath, flag.Args())
}

func findNotePath() string {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot get current directory: %s\n", err)
		os.Exit(1)
	}
	separator := string(filepath.Separator)
	for dir := wd; dir != "." && dir != separator; dir = filepath.Dir(dir) {
		notePath := filepath.Join(dir, NoteDir)
		fi, err := os.Stat(notePath)
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "cannot stat %s: %s. ignoring\n", notePath, err)
			}
			continue
		}
		if fi.IsDir() {
			rel, err := filepath.Rel(wd, notePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot get relative path to %s: %s. ignoring\n", notePath, err)
				continue
			}
			return rel
		}
	}
	return NoteDir
}

func run(notePath string, args []string) {
	if len(args) > 0 {
		command := args[0]
		if fn, ok := handlers[command]; ok {
			fn(notePath, args[1:])
		} else {
			fmt.Fprintf(os.Stderr, "unknown command %q\n", command)
			flag.Usage()
			os.Exit(2)
		}
	} else {
		fmt.Fprintln(os.Stderr, "no command specified")
		flag.Usage()
		os.Exit(2)
	}
}
