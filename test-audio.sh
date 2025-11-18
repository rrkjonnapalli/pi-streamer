#!/bin/bash

echo "=== Audio Device Test ==="
echo ""

echo "1. Listing audio capture devices:"
arecord -l
echo ""

echo "2. Testing audio capture for 3 seconds..."
echo "   (Speak into your microphone now!)"
arecord -D hw:0,0 -f S16_LE -c 1 -r 48000 -d 3 test-capture.wav
echo ""

echo "3. Checking captured audio file:"
if [ -f test-capture.wav ]; then
    SIZE=$(stat -f%z test-capture.wav 2>/dev/null || stat -c%s test-capture.wav 2>/dev/null)
    echo "   File size: $SIZE bytes"
    echo "   Expected: ~288KB for 3 seconds"
    
    if [ $SIZE -lt 10000 ]; then
        echo "   ⚠️  WARNING: File is too small - microphone may not be working!"
    else
        echo "   ✓ File size looks good"
    fi
    echo ""
    
    echo "4. Playing back captured audio..."
    echo "   (You should hear what you just said)"
    aplay test-capture.wav
    echo ""
    
    echo "5. Cleaning up..."
    rm -f test-capture.wav
    echo "   Test complete!"
else
    echo "   ❌ ERROR: No audio file created - microphone not working!"
fi

echo ""
echo "=== Volume Levels ==="
echo "Checking mixer settings..."
amixer scontrols
echo ""
echo "To adjust volume, use: amixer set <control> <percentage>%"
echo "Example: amixer set Mic 80%"
