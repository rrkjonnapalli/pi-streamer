package main

import (
	"log"
	"time"

	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
)

func StartPublisher(cfg Config) error {
	disconnected := make(chan struct{})

	room, err := lksdk.ConnectToRoom(cfg.LivekitURL, lksdk.ConnectInfo{
		APIKey:              cfg.APIKey,
		APISecret:           cfg.APISecret,
		RoomName:            cfg.RoomName,
		ParticipantIdentity: cfg.PublisherName,
	}, &lksdk.RoomCallback{
		ParticipantCallback: lksdk.ParticipantCallback{
			OnTrackPublished: func(publication *lksdk.RemoteTrackPublication, rp *lksdk.RemoteParticipant) {
				log.Println("track published")
			},
		},
		OnDisconnected: func() {
			log.Println("room disconnected")
			close(disconnected)
		},
	})
	if err != nil {
		return err
	}
	defer room.Disconnect()

	log.Println("connected to room:", cfg.RoomName)

	track, err := webrtc.NewTrackLocalStaticSample(
		webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus},
		"mic",
		"audio",
	)
	if err != nil {
		return err
	}

	_, err = room.LocalParticipant.PublishTrack(track, &lksdk.TrackPublicationOptions{
		Name:   "mic",
		Source: livekit.TrackSource_MICROPHONE,
	})
	if err != nil {
		return err
	}

	log.Println("track published, starting audio stream")

	capture, err := NewAudioCapture(cfg)
	if err != nil {
		return err
	}
	defer capture.Stop()

	if err := capture.Start(); err != nil {
		return err
	}

	sequencer := rtp.NewRandomSequencer()
	frameCount := uint64(0)
	lastLog := time.Now()

	for {
		select {
		case <-disconnected:
			return nil
		default:
		}

		opusData, err := capture.ReadOpusFrame()
		if err != nil {
			log.Println("audio read error:", err)
			return err
		}

		if err := track.WriteSample(media.Sample{
			Data:     opusData,
			Duration: time.Millisecond * 20,
		}); err != nil {
			log.Println("track write error:", err)
			return err
		}

		frameCount++
		sequencer.NextSequenceNumber()

		if time.Since(lastLog) > 5*time.Minute {
			log.Printf("streaming active: %d frames sent (%.1f hours)", frameCount, float64(frameCount)*0.02/3600)
			lastLog = time.Now()
		}
	}
}
