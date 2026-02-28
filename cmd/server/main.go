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

	"coffee-tracker-backend/internal/server"
)

func main() {
	// Create a root context that cancels on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize the server
	srv, err := server.NewServer()
	if err != nil {
		log.Fatalf("❌ Failed to initialize server: %v", err)
	}

	// Start the server in a separate goroutine
	go func() {
		if err := srv.Start(); err != nil {
			srv.Logger.Fatalf("Server terminated with error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	srv.Logger.Println("🛑 Shutdown signal received, stopping server gracefully...")

	// Create a timeout context for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Perform graceful shutdown
	if err := srv.Shutdown(shutdownCtx); err != nil {
		srv.Logger.Printf("⚠️ Error during shutdown: %v", err)
	} else {
		srv.Logger.Println("✅ Server stopped cleanly")
	}
}
