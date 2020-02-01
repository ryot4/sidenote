package main

import (
	"flag"
	"fmt"
	"os"
)

func runInit(noteDir string, args []string) {
	var linkDir string

	initFlag := flag.NewFlagSet("init", flag.ExitOnError)
	initFlag.StringVar(&linkDir, "l", "", fmt.Sprintf("make %s a symlink to the directory", noteDir))
	initFlag.Parse(args)

	var err error
	if linkDir == "" {
		_, err = InitDirectory(noteDir)
	} else {
		_, err = InitDirectoryLink(linkDir, noteDir)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "init failed: %s\n", err)
		os.Exit(1)
	}
}
