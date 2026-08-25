package main

import (
	"log"

	"adamndegwa/internal/server"
)

func main() {
	log.Fatal(server.RunFromEnv())
}
