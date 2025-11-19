# Pi Streamer

LiveKit audio streamer for Raspberry Pi (also works on macOS/Linux)

## Installation

### Raspberry Pi / Ubuntu:
```bash
sudo apt update
sudo apt install -y portaudio19-dev libopus-dev pkg-config golang-go
```

### macOS:
```bash
brew install portaudio opus pkg-config go
```

## Setup

### 1. Configure Environment
Copy the example environment file:
```bash
cp .env.example .env
```

Edit `.env` with your LiveKit credentials:
```bash
nano .env
```

Required variables:
```env
LIVEKIT_URL=wss://your-livekit-server.com
LIVEKIT_API_KEY=your-api-key
LIVEKIT_API_SECRET=your-api-secret
ROOM_NAME=your-room-name
PI_PUBLISHER_NAME=pi-audio-publisher
OPUS_BITRATE=16000
```

**Note:** Uses system default audio input device.

### 2. Build and Run
```bash
go mod download
go build -o dist/streamer
./dist/streamer
```

## Running as Service (24/7)

### Install Service

```bash
# Update paths in service/pistreamer.service if needed
sudo cp service/pistreamer.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable pistreamer
sudo systemctl start pistreamer

# Check status and logs
sudo systemctl status pistreamer
sudo journalctl -u pistreamer -f
```

## Uninstall

```bash
sudo systemctl stop pistreamer
sudo systemctl disable pistreamer
sudo rm /etc/systemd/system/pistreamer.service
sudo systemctl daemon-reload
```
