package main

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	LivekitURL    string
	APIKey        string
	APISecret     string
	RoomName      string
	PublisherName string
	ALSADevice    string
	OpusBitrate   int
}

func detectALSADevice() string {
	cmd := exec.Command("arecord", "-l")
	output, err := cmd.Output()
	if err != nil {
		log.Println("failed to detect ALSA device, using default")
		return "default"
	}

	re := regexp.MustCompile(`card (\d+):.*device (\d+):`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) >= 3 {
		device := "hw:" + matches[1] + "," + matches[2]
		log.Println("auto-detected ALSA device:", device)
		return device
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "card") {
			log.Println("ALSA devices found but couldn't parse, using default")
			return "default"
		}
	}

	log.Println("no ALSA device found, using default")
	return "default"
}

func LoadConfig() Config {
	bitrate := 16000
	if br := os.Getenv("OPUS_BITRATE"); br != "" {
		if parsed, err := strconv.Atoi(br); err == nil {
			bitrate = parsed
		}
	}

	alsaDevice := os.Getenv("DEVICE_ALSA_HW")
	if alsaDevice == "" {
		alsaDevice = detectALSADevice()
	}

	return Config{
		LivekitURL:    os.Getenv("LIVEKIT_URL"),
		APIKey:        os.Getenv("LIVEKIT_API_KEY"),
		APISecret:     os.Getenv("LIVEKIT_API_SECRET"),
		RoomName:      os.Getenv("ROOM_NAME"),
		PublisherName: os.Getenv("PI_PUBLISHER_NAME"),
		ALSADevice:    alsaDevice,
		OpusBitrate:   bitrate,
	}
}
