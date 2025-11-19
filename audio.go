package main

import (
	"fmt"
	"log"

	"github.com/gordonklaus/portaudio"
	"gopkg.in/hraban/opus.v2"
)

const (
	SampleRate = 48000
	Channels   = 1
	FrameSize  = 960
	MaxPacket  = 4000
)

type AudioCapture struct {
	stream  *portaudio.Stream
	encoder *opus.Encoder
	buffer  []int16
	cfg     Config
}

func NewAudioCapture(cfg Config) (*AudioCapture, error) {
	encoder, err := opus.NewEncoder(SampleRate, Channels, opus.AppVoIP)
	if err != nil {
		return nil, err
	}
	encoder.SetBitrate(cfg.OpusBitrate)
	encoder.SetDTX(false)
	encoder.SetPacketLossPerc(0)

	if err := portaudio.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize portaudio: %w", err)
	}

	log.Printf("Opus encoder: bitrate=%d, DTX=disabled, VBR=enabled", cfg.OpusBitrate)

	return &AudioCapture{
		encoder: encoder,
		buffer:  make([]int16, FrameSize),
		cfg:     cfg,
	}, nil
}

func (ac *AudioCapture) Start() error {
	var err error
	ac.stream, err = portaudio.OpenDefaultStream(Channels, 0, SampleRate, FrameSize, ac.buffer)
	if err != nil {
		return fmt.Errorf("failed to open audio stream: %w", err)
	}

	if err := ac.stream.Start(); err != nil {
		return fmt.Errorf("failed to start audio stream: %w", err)
	}

	log.Println("audio capture started (PortAudio)")
	return nil
}

func (ac *AudioCapture) ReadOpusFrame() ([]byte, error) {
	if err := ac.stream.Read(); err != nil {
		return nil, fmt.Errorf("audio read failed: %w", err)
	}

	opusData := make([]byte, MaxPacket)
	n, err := ac.encoder.Encode(ac.buffer, opusData)
	if err != nil {
		return nil, fmt.Errorf("opus encode failed: %w", err)
	}

	if n < 20 {
		log.Printf("Very small Opus packet: %d bytes (buffer has data: first=%d, max=%d)",
			n, ac.buffer[0], maxAbs(ac.buffer))
	}

	return opusData[:n], nil
}

func maxAbs(samples []int16) int16 {
	var max int16
	for _, s := range samples {
		if s < 0 {
			s = -s
		}
		if s > max {
			max = s
		}
	}
	return max
}

func (ac *AudioCapture) GetBufferRMS() float64 {
	var sum float64
	for _, sample := range ac.buffer {
		sum += float64(sample) * float64(sample)
	}
	return sum / float64(len(ac.buffer))
}

func (ac *AudioCapture) GetPCMLevel(pcm []int16) float64 {
	var sum float64
	for _, sample := range pcm {
		sum += float64(sample) * float64(sample)
	}
	return sum / float64(len(pcm))
}

func (ac *AudioCapture) Stop() {
	if ac.stream != nil {
		ac.stream.Stop()
		ac.stream.Close()
	}
	portaudio.Terminate()
}
