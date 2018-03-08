package main

import (
	"log"
	"time"

	"mcquay.me/hw"
)

func main() {
	for {
		log.Printf("hwl@%+v", hw.Version)
		time.Sleep(1 * time.Second)
	}
}
