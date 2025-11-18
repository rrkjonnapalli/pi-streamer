package main

import (
	"io"
	"log"
	"os/exec"

	"gopkg.in/hraban/opus.v2"
)

const (
	SampleRate = 48000
	Channels   = 1
	FrameSize  = 960
	MaxPacket  = 4000
)

type AudioCapture struct {
	cmd     *exec.Cmd
	stdout  io.ReadCloser
	encoder *opus.Encoder
	cfg     Config
}

func NewAudioCapture(cfg Config) (*AudioCapture, error) {
	encoder, err := opus.NewEncoder(SampleRate, Channels, opus.AppVoIP)
	if err != nil {
		return nil, err
	}
	encoder.SetBitrate(cfg.OpusBitrate)

	return &AudioCapture{
		encoder: encoder,
		cfg:     cfg,
	}, nil
}

func (ac *AudioCapture) Start() error {
	ac.cmd = exec.Command("arecord",
		"-D", ac.cfg.ALSADevice,
		"-f", "S16_LE",
		"-c", "1",
		"-r", "48000",
		"-t", "raw",
	)

	stdout, err := ac.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	ac.stdout = stdout

	stderr, err := ac.cmd.StderrPipe()
	if err != nil {
		return err
	}

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				log.Printf("arecord: %s", string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	if err := ac.cmd.Start(); err != nil {
		return err
	}

	log.Println("audio capture started")
	return nil
}

func (ac *AudioCapture) ReadOpusFrame() ([]byte, error) {
	pcm := make([]int16, FrameSize)
	pcmBytes := make([]byte, FrameSize*2)

	_, err := io.ReadFull(ac.stdout, pcmBytes)
	if err != nil {
		return nil, err
	}

	for i := 0; i < FrameSize; i++ {
		pcm[i] = int16(pcmBytes[i*2]) | int16(pcmBytes[i*2+1])<<8
	}

	opusData := make([]byte, MaxPacket)
	n, err := ac.encoder.Encode(pcm, opusData)
	if err != nil {
		return nil, err
	}

	return opusData[:n], nil
}

func (ac *AudioCapture) GetPCMLevel(pcm []int16) float64 {
	var sum float64
	for _, sample := range pcm {
		sum += float64(sample) * float64(sample)
	}
	return sum / float64(len(pcm))
}

func (ac *AudioCapture) Stop() {
	if ac.cmd != nil && ac.cmd.Process != nil {
		ac.cmd.Process.Kill()
	}
	if ac.stdout != nil {
		ac.stdout.Close()
	}
}
