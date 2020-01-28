package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"
)

func runList(args []string) {
	var longFormat, sortByMtime bool

	listFlag := flag.NewFlagSet("list", flag.ExitOnError)
	listFlag.BoolVar(&longFormat, "l", false, "long format")
	listFlag.BoolVar(&sortByMtime, "t", false, "sort by modification time")
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

	listItems(files, longFormat, sortByMtime)
	for i, fi := range dirs {
		if len(listPaths) > 1 {
			fmt.Printf("\n%s:\n", fi.Name())
		}
		listDir(dir, dirPaths[i], longFormat, sortByMtime)
	}
}

func listDir(dir *Directory, path string, longFormat, sortByMtime bool) {
	children, err := dir.Readdir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read %s: %s\n", path, err)
	}
	listItems(children, longFormat, sortByMtime)
}

func listItems(items []os.FileInfo, longFormat, sortByMtime bool) {
	if sortByMtime {
		sort.Slice(items, func(i, j int) bool {
			return items[i].ModTime().After(items[j].ModTime())
		})
	} else {
		sort.Slice(items, func(i, j int) bool {
			return items[i].Name() < items[j].Name()
		})
	}

	for _, fi := range items {
		printFile(fi, longFormat)
	}
}

func printFile(fi os.FileInfo, longFormat bool) {
	name := fi.Name()
	if fi.IsDir() {
		name += "/"
	}
	if longFormat {
		fmt.Printf("%s %s\n", formatTime(fi.ModTime()), name)
	} else {
		fmt.Println(name)
	}
}

func formatTime(t time.Time) string {
	if t.Year() == time.Now().Year() {
		return t.Format("Jan _2 15:04")
	} else {
		return t.Format("Jan _2  2006")
	}
}
