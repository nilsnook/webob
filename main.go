package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// handling SIGTERM and SIGINT signals
	// creating a channel to listen for these signals
	sigChan := make(chan os.Signal, 1)
	// setup notifications for these signals
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// separate goroutine to listen for these signals once notified
	go func() {
		select {
		case _ = <-sigChan:
			log.Println("Got SIGINT/SIGTERM, exiting.")
			cancel()
			os.Exit(1)
		case <-ctx.Done():
			log.Fatalln("Done.")
		}
	}()

	defer func() {
		signal.Stop(sigChan)
		cancel()
	}()

	c := &config{}
	// passing in config and
	// setting default log output to STDOUT
	log.Fatalln(run(ctx, c, os.Stdout))
}

func run(ctx context.Context, c *config, out io.Writer) error {
	// init config
	err := c.init()
	if err != nil {
		return err
	}
	// set output log
	log.SetOutput(out)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(c.Tick):
			// hit the given URL
			resp, err := http.Get(c.Url)
			if err != nil {
				return err
			}

			// parse response meta data
			// and match against user defined configurations.
			//
			// status code
			if resp.StatusCode != c.StatusCode {
				logError("Server code mismatch", strconv.Itoa(resp.StatusCode), strconv.Itoa(c.StatusCode))
			}

			// server
			if s := resp.Header.Get("server"); s != c.Server {
				logError("Server header mismatch", s, c.Server)
			}

			// content type
			if ct := resp.Header.Get("content-type"); ct != c.ContentType {
				logError("Content-Type header mismatch", ct, c.ContentType)
			}

			// user agent
			if ua := resp.Header.Get("user-agent"); ua != c.UserAgent {
				logError("User-Agent header mismatch", ua, c.UserAgent)
			}
		}
	}
}
