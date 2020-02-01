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

type HandlerFunc func(noteDir string, args []string)

var handlers = map[string]HandlerFunc{
	"init": runInit,
	"edit": runEdit,
	"ls":   runLs,
}

func main() {
	noteDir := flag.String("d", findNoteDir(), "path to the directory for notes")
	flag.Parse()
	run(*noteDir, flag.Args())
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

func run(noteDir string, args []string) {
	if len(args) > 0 {
		command := args[0]
		if fn, ok := handlers[command]; ok {
			fn(noteDir, args[1:])
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
