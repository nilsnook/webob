package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// set log output to STDOUT explicitly
	log.SetOutput(os.Stdout)

	// context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// config
	c := &config{}
	// init config
	err := c.init()
	if err != nil {
		log.Fatalln(err)
	}

	// handling SIGTERM and SIGINT termination signals and
	// SIGHUP to reload configuration
	//
	// creating a channel to listen for these signals
	sigChan := make(chan os.Signal, 1)
	// setup notifications for these signals
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	// separate goroutine to listen for these signals once notified
	go func() {
		for {
			select {
			case s := <-sigChan:
				switch s {
				case syscall.SIGINT, syscall.SIGTERM:
					log.Println("Got SIGINT/SIGTERM, exiting.")
					cancel()
					os.Exit(1)
				case syscall.SIGHUP:
					log.Println("Got SIGHUP, reloading with new config.")
					// reinitialize config
					nc := &config{}
					err := nc.initFromConfigFile()
					if err != nil {
						log.Printf("ERROR: Error initializing with new config settings.\n")
						log.Printf("\t%s\n", err)
						log.Printf("\tResuming with old config...\n")
					} else {
						// set new configuration
						c = nc
					}
				}
			case <-ctx.Done():
				log.Fatalln("Done.")
			}
		}
	}()

	defer func() {
		signal.Stop(sigChan)
		cancel()
	}()

	// run with config
	log.Fatalln(run(ctx, c))
}

func run(ctx context.Context, c *config) error {
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
