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
		fmt.Fprintf(output, "Usage: %s %s [-l address[:port]]\n", os.Args[0], c.Name())
		fmt.Fprintf(output, "\n%s.\n", c.Description())
		fmt.Fprintln(output, "\noptions:")
		c.flag.PrintDefaults()
	}
	c.flag.StringVar(&c.listenAddress, "l", "0.0.0.0:8000", "Address to listen")
	c.flag.Parse(args)
}

func (c *ServeCommand) Run(args []string, options *Options) error {
	c.setup(args, options)

	dir, err := checkDirectory(options.noteDir)
	if err != nil {
		return err
	}

	srv := NewServer(c.listenAddress, dir.path)

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

	log.Println("Listening on", c.listenAddress)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err.Error())
	}

	<-idleConnsClosed
	return nil
}
