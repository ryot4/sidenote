package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type ServeCommand struct {
	flag *flag.FlagSet

	contentType   string
	listenAddress string
}

func (c *ServeCommand) Name() string {
	return "serve"
}

func (c *ServeCommand) Description() string {
	return "Serve notes over HTTP"
}

func (c *ServeCommand) setup(args []string, options *Options) {
	c.flag = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	c.flag.Usage = func() {
		output := c.flag.Output()
		fmt.Fprintf(output, "Usage: %s %s [-l address[:port]] [-t content-type] [subdir]\n",
			os.Args[0], c.Name())
		fmt.Fprintf(output, "\n%s.\n", c.Description())
		fmt.Fprintln(output, "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.StringVar(&c.listenAddress, "l", "0.0.0.0:8000", "Address to listen on")
	c.flag.StringVar(&c.contentType, "t", "", "Specify Content-Type of notes")
	c.flag.Parse(args)
}

func (c *ServeCommand) Run(args []string, options *Options) error {
	c.setup(args, options)

	dir, err := checkDirectory(options.noteDir)
	if err != nil {
		return err
	}

	if c.flag.NArg() > 1 {
		return NewSyntaxError("too many arguments")
	}
	subdir := c.flag.Arg(0)
	if subdir == "" {
		subdir = "/"
	}

	return c.runServer(dir, subdir)
}

func (c *ServeCommand) runServer(dir *Directory, subdir string) error {
	root, err := dir.JoinPath(subdir)
	if err != nil {
		return err
	}

	if fi, err := os.Stat(root); err != nil {
		return err
	} else if !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", subdir)
	}

	srv := NewServer(c.listenAddress, root, c.contentType)

	idleConnsClosed := make(chan struct{})
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		signal.Notify(sigCh, syscall.SIGTERM)
		<-sigCh

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Println("Failed to shutdown the server gracefully:", err.Error())
		}
		close(idleConnsClosed)
	}()

	log.Printf("listening on %s (root directory: %s)\n", c.listenAddress, root)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Print(err.Error())
		return err
	}

	<-idleConnsClosed
	return nil
}
