package main

import (
	"distributed.systems.labs/worker/internal"
	"log"
)

func main() {
	log.Printf("starting worker")
	internal.Main()
}
