// file: cmd/server/main.go
// Main entry point for the Coffee Tracker backend API server.
package main

import (
	"log"
)

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	if err := server.Start(); err != nil {
		server.logger.Fatalf("Server terminated with error: %v", err)
	}
}