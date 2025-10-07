// file: cmd/server/main.go
// Main entry point for the Coffee Tracker backend API server.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Create a root context that cancels on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize the server
	server, err := NewServer()
	if err != nil {
		log.Fatalf("❌ Failed to initialize server: %v", err)
	}

	// Start the server in a separate goroutine
	go func() {
		if err := server.Start(); err != nil {
			server.logger.Fatalf("Server terminated with error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	server.logger.Println("🛑 Shutdown signal received, stopping server gracefully...")

	// Create a timeout context for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Perform graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		server.logger.Printf("⚠️ Error during shutdown: %v", err)
	} else {
		server.logger.Println("✅ Server stopped cleanly")
	}
}
