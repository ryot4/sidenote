package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	NotePath = ".notes"
)

type HandlerFunc func(args []string)

var handlers = map[string]HandlerFunc{
	"init": runInit,
	"edit": runEdit,
	"ls":   runLs,
}

func main() {
	flag.Parse()
	run(flag.Args())
}

func run(args []string) {
	if len(args) > 0 {
		command := args[0]
		if fn, ok := handlers[command]; ok {
			fn(args[1:])
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
