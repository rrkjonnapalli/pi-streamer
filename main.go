package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	cfg := LoadConfig()

	if cfg.LivekitURL == "" || cfg.APIKey == "" || cfg.APISecret == "" || cfg.RoomName == "" {
		log.Fatal("missing required environment variables")
	}

	log.Printf("starting pi-streamer: room=%s, bitrate=%d", cfg.RoomName, cfg.OpusBitrate)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("shutting down...")
		os.Exit(0)
	}()

	retryCount := 0
	for {
		err := StartPublisher(cfg)
		if err != nil {
			log.Printf("publisher exited with error (retry %d): %v", retryCount, err)
			retryCount++
		} else {
			log.Println("publisher exited cleanly")
			retryCount = 0
		}

		backoff := time.Duration(min(retryCount, 10)) * time.Second
		if backoff < 3*time.Second {
			backoff = 3 * time.Second
		}

		log.Printf("reconnecting in %v...", backoff)
		time.Sleep(backoff)
	}
}
