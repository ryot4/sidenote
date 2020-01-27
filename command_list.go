package main

import (
	"flag"
	"fmt"
	"os"
)

func runList(args []string) {
	var longFormat bool

	listFlag := flag.NewFlagSet("list", flag.ExitOnError)
	listFlag.BoolVar(&longFormat, "l", false, "long format")
	listFlag.Parse(args)

	dir, err := OpenDirectory(NotePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	listPaths := listFlag.Args()
	if len(listPaths) == 0 {
		listPaths = append(listPaths, "/")
	}

	var files, dirs []os.FileInfo
	var dirPaths []string
	for _, path := range listPaths {
		fi, err := dir.Stat(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot stat %s: %s\n", path, err)
		}
		if fi.IsDir() {
			dirs = append(dirs, fi)
			dirPaths = append(dirPaths, path)
		} else {
			files = append(files, fi)
		}
	}

	for _, fi := range files {
		printFile(fi, longFormat)
	}
	for i, fi := range dirs {
		if len(listPaths) > 1 {
			fmt.Printf("\n%s:\n", fi.Name())
		}
		listDir(dir, dirPaths[i], longFormat)
	}
}

func printFile(fi os.FileInfo, longFormat bool) {
	name := fi.Name()
	if fi.IsDir() {
		name += "/"
	}
	if longFormat {
		fmt.Printf("%s %s\n",
			fi.ModTime().Format("Jan _2 15:04"),
			name)
	} else {
		fmt.Println(name)
	}
}

func listDir(dir *Directory, path string, longFormat bool) {
	children, err := dir.Readdir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read %s: %s\n", path, err)
	}
	for _, child := range children {
		printFile(child, longFormat)
	}
}
