package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()

	c := &config{}
	log.Fatalln(run(ctx, c))
}

func run(ctx context.Context, c *config) error {
	err := c.init()
	if err != nil {
		return err
	}

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
