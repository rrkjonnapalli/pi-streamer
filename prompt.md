Build a Raspberry Pi audio streamer that runs 24/7 and publishes a continuous Opus WebRTC audio track to a LiveKit server. This must run on Ubuntu for Raspberry Pi and must require zero code edits after initial setup. All configuration comes only from a .env file.


Requirements
1.	Project structure:
```
pi-streamer/
  main.go
  audio.go
  webrtc.go
  config.go
  go.mod
  .env.example
  .env (ignored in git)
  run.sh
  service/pistreamer.service
```

2.	.env.example must contain these variables:

```
LIVEKIT_URL=
LIVEKIT_API_KEY=
LIVEKIT_API_SECRET=
ROOM_NAME=
PI_PUBLISHER_NAME=
DEVICE_ALSA_HW=
OPUS_BITRATE=
```

3.	On start:

•	Load .env
•	Connect to LiveKit using URL + API key + secret
•	Join the room from .env
•	Create a single Opus audio track named mic
•	Start continuous streaming

4.	Audio capture:

•	Capture 20ms mono PCM frames from ALSA using arecord
•	Sample rate: 48000 Hz
•	Format: S16_LE
•	Encode frames as Opus (Pion handles this automatically when writing samples)
•	If capture fails, reconnect automatically

5.	Stream behavior:

•	WebRTC connection must auto-recover if disconnected
•	Always send Opus frames (real audio or silence)
•	No additional logic, VAD, or talkback yet

6.	run.sh:

•	Loads .env
•	Starts the binary

7.	Systemd service:

•	Runs the streamer at boot
•	Restarts automatically
•	Path: /etc/systemd/system/pistreamer.service

8.	Result:

•	After editing .env, user runs:
```bash
go build -o dist/streamer
./run.sh
```
It must immediately connect to LiveKit and stream audio 24/7.

9.	Do not add features not listed. Keep the implementation minimal and deterministic.
