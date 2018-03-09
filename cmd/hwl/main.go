package main

import (
	"log"
	"time"

	"mcquay.me/hw"
)

func main() {
	for {
		log.Printf("hwl version=%q, git=%q", hw.Version, hw.Git)
		time.Sleep(1 * time.Second)
	}
}
