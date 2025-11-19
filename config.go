package main

import (
	"os"
	"strconv"
)

type Config struct {
	LivekitURL    string
	APIKey        string
	APISecret     string
	RoomName      string
	PublisherName string
	OpusBitrate   int
}

func LoadConfig() Config {
	bitrate := 16000
	if br := os.Getenv("OPUS_BITRATE"); br != "" {
		if parsed, err := strconv.Atoi(br); err == nil {
			bitrate = parsed
		}
	}

	return Config{
		LivekitURL:    os.Getenv("LIVEKIT_URL"),
		APIKey:        os.Getenv("LIVEKIT_API_KEY"),
		APISecret:     os.Getenv("LIVEKIT_API_SECRET"),
		RoomName:      os.Getenv("ROOM_NAME"),
		PublisherName: os.Getenv("PI_PUBLISHER_NAME"),
		OpusBitrate:   bitrate,
	}
}
