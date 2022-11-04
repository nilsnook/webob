package main

import (
	"fmt"
	"log"
)

func logError(msg string, received string, expected string) {
	info := fmt.Sprintf("%s\n\tReceived: %s\n\tExpected: %s", msg, received, expected)
	log.Println(info)
}
