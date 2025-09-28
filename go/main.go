package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"music-app/src/database"
	"music-app/src/handlers"
	"music-app/src/middleware"
	"music-app/src/otel"

	_ "modernc.org/sqlite"
)

func main() {
	// Initialize OpenTelemetry
	shutdown := otel.InitOTel("music-app", "1.0.0")
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			log.Printf("Error shutting down OpenTelemetry: %v", err)
		}
	}()

	// Initialize database
	db, err := sql.Open("sqlite", "./music.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create instrumented database wrapper
	instrumentedDB := database.NewInstrumentedDB(db)

	// Create tables
	if err := database.CreateTables(db); err != nil {
		log.Fatal(err)
	}

	// Insert sample data
	if err := database.InsertSampleData(db); err != nil {
		log.Fatal(err)
	}

	// Create handler with instrumented DB
	songHandler := handlers.NewSongHandler(instrumentedDB.DB) // Use the underlying DB for now

	// Set up routes with instrumentation middleware
	http.HandleFunc("/test", middleware.HTTPTelemetryMiddleware(songHandler.GetSongWithLyrics, "get_song_with_lyrics"))

	// Set up graceful shutdown
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		fmt.Println("Server starting on port 8081...")
		serverErrors <- http.ListenAndServe(":8081", nil)
	}()

	// Wait for interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)
	case <-interrupt:
		fmt.Println("Shutting down server...")
	}
}
