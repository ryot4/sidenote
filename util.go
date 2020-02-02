package main

import (
	"fmt"
	"os"
)

func exitWithSyntaxError(message string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", message)
	os.Exit(2)
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	os.Exit(1)
}
