# Pi Streamer - Setup and Run Guide

## System Requirements
- Raspberry Pi 5 running Ubuntu
- Audio input device (USB microphone or audio interface)
- Network connection to LiveKit server

## Dependencies Installation

### 1. Install System Packages
```bash
sudo apt update
sudo apt install -y alsa-utils libasound2-dev libopus-dev pkg-config golang-go
```

### 2. Verify Go Installation
```bash
go version
```
Should show Go 1.22 or higher. If not, install from [golang.org](https://golang.org/dl/).

### 3. Test Audio Device
List available audio devices:
```bash
arecord -l
```
Note your device hardware ID (e.g., `hw:1,0`).

Test recording:
```bash
arecord -D hw:1,0 -f S16_LE -c 1 -r 48000 -d 5 test.wav
aplay test.wav
```

## Project Setup

### 1. Clone/Copy Project Files
```bash
cd ~
# Your project should be in ~/pi-streamer
cd pi-streamer
```

### 2. Configure Environment
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
DEVICE_ALSA_HW=hw:1,0
OPUS_BITRATE=16000
```

**Note:** `DEVICE_ALSA_HW` is optional. If not provided, the system will automatically detect the first available audio capture device using `arecord -l`.

### 3. Download Dependencies
```bash
go mod download
```

### 4. Build the Binary
```bash
o build -o dist/streamer
```

### 5. Make Run Script Executable
```bash
chmod +x run.sh
```

## Running the Streamer

### Manual Run (Testing)
```bash
./run.sh
```

You should see:
```
starting pi-streamer: room=your-room-name, device=hw:1,0, bitrate=16000
connected to room: your-room-name
track published, starting audio stream
audio capture started
```

Stop with `Ctrl+C`.

### Run as System Service (24/7)

1. Edit the service file if needed:
```bash
nano service/pistreamer.service
```

Update `WorkingDirectory` and `ExecStart` paths if your installation is not in `/home/pi/pi-streamer`.

2. Copy service file:
```bash
sudo cp service/pistreamer.service /etc/systemd/system/
```

3. Enable and start service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable pistreamer
sudo systemctl start pistreamer
```

4. Check status:
```bash
sudo systemctl status pistreamer
```

5. View logs:
```bash
sudo journalctl -u pistreamer -f
```

## Troubleshooting

### Audio Device Not Found
The system auto-detects your audio device if `DEVICE_ALSA_HW` is not set. To manually specify:
```bash
arecord -l
```
Update `DEVICE_ALSA_HW` in `.env` with correct device (e.g., `hw:1,0`).

### Permission Denied for Audio
Add user to audio group:
```bash
sudo usermod -a -G audio $USER
```
Reboot required.

### Connection Issues
- Verify `LIVEKIT_URL` starts with `wss://`
- Check API key and secret are correct
- Test network: `ping your-livekit-server.com`

### Service Won't Start
Check logs:
```bash
sudo journalctl -u pistreamer -n 50
```

### High CPU Usage
Lower bitrate in `.env`:
```env
OPUS_BITRATE=12000
```

## Notes

### Behavior
- **Silence handling:** Opus codec efficiently encodes silence (small packets ~20-50 bytes)
- **Speech handling:** When someone speaks, packets grow to 100-200+ bytes, sent in real-time
- **Continuous operation:** Designed for 24/7/365 operation with automatic recovery
- Audio is encoded as Opus at 48kHz mono, each frame is 20ms (960 samples)
- Default bitrate is 16kbps (configurable)
- No local audio storage - streams directly to LiveKit

### Reliability Features
- **Auto-reconnection:** Exponential backoff (3-30s) if connection drops
- **Connection monitoring:** Detects LiveKit disconnects and reconnects automatically
- **Heartbeat logging:** Logs streaming statistics every 5 minutes
- **ALSA monitoring:** Captures and logs any audio device warnings/errors
- **Clean shutdown:** Handles SIGINT/SIGTERM signals gracefully

### System Requirements
- Service runs as the user specified in `pistreamer.service` (default: `pi`)
- Requires persistent network connection to LiveKit server
- Minimal CPU usage (~5-10% on Raspberry Pi 5)
- Memory footprint: ~30-50MB

## Updating Configuration

After changing `.env`:

**If running manually:**
Stop with `Ctrl+C` and restart: `./run.sh`

**If running as service:**
```bash
sudo systemctl restart pistreamer
```

## Production Considerations

### For 24/7 Operation:
1. **Monitor logs regularly:**
   ```bash
   sudo journalctl -u pistreamer -f
   ```

2. **Set up log rotation** (optional but recommended):
   ```bash
   sudo nano /etc/systemd/journald.conf
   # Set: SystemMaxUse=500M
   sudo systemctl restart systemd-journald
   ```

3. **Network stability:**
   - Ensure stable network connection
   - Consider static IP or DHCP reservation
   - Monitor internet connectivity

4. **Power management:**
   - Disable sleep/suspend on Raspberry Pi
   - Consider UPS for power backup

5. **Disk space:**
   - Monitor available disk space (logs can grow)
   - The streamer itself doesn't write files, but system logs do

### Known Limitations:
- Network outages longer than a few minutes will accumulate reconnection attempts
- ALSA device must remain available (don't unplug USB mic while running)
- LiveKit room must exist (create room before connecting publisher)

## Uninstall Service

```bash
sudo systemctl stop pistreamer
sudo systemctl disable pistreamer
sudo rm /etc/systemd/system/pistreamer.service
sudo systemctl daemon-reload
```
