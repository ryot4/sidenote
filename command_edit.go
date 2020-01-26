package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func runEdit(args []string) {
	var editor string

	editFlag := flag.NewFlagSet("edit", flag.ExitOnError)
	editFlag.StringVar(&editor, "e", os.Getenv("EDITOR"), "editor to use")
	editFlag.Parse(args)

	dir, err := OpenDirectory(NotePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(editFlag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "no filename specified")
		os.Exit(1)
	}
	path, err := dir.FilePath(editFlag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid path: %s\n", err)
		os.Exit(2)
	}
	runEditor(editor, path)
}

func runEditor(editor, path string) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "cannot create %s: %v\n", dir, err)
		os.Exit(1)
	}

	fi, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "cannot stat %s: %v\n", path, err)
			os.Exit(1)
		}
	} else if fi.IsDir() {
		fmt.Fprintf(os.Stderr, "directory exists: %s\n", path)
		os.Exit(1)
	}

	editorCmd := exec.Command(editor, path)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr
	err = editorCmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	os.Exit(editorCmd.ProcessState.ExitCode())
}
