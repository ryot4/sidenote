package main

import (
	"flag"
	"fmt"
	"os"
)

func runInit(args []string) {
	var linkDir string

	initFlag := flag.NewFlagSet("init", flag.ExitOnError)
	initFlag.StringVar(&linkDir, "l", "", fmt.Sprintf("make %s a symlink to the directory", NotePath))
	initFlag.Parse(args)

	var err error
	if linkDir == "" {
		_, err = InitDirectory(NotePath)
	} else {
		_, err = InitDirectoryLink(linkDir, NotePath)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "init failed: %s\n", err)
		os.Exit(1)
	}
}
