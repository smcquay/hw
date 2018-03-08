package main

import (
	"log"
	"time"
)

const version = "v0.1.1"

func main() {
	for {
		log.Printf("hwl@%+v", version)
		time.Sleep(1 * time.Second)
	}
}
