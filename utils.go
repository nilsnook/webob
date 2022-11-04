package main

import (
	"fmt"
	"log"
)

func logError(msg string, received string, expected string) {
	out := fmt.Sprintf("%s\nReceived: %s\nExpected: %s", msg, received, expected)
	log.Println(out)
}
