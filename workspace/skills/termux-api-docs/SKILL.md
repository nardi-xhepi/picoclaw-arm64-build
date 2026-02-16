---
name: termux-api-docs
description: Reference manual for Termux API commands. Use this to learn how to control Android hardware (camera, battery, clipboard, etc.) via the command line.
---

# Termux API Reference

This skill provides documentation for the Termux API. You can use the `exec` tool to run these commands directly.

> **Requirement**: The user must have the `Termux:API` app installed and the `termux-api` package (`pkg install termux-api`).

## Hardware & Sensors

### Camera
**Take a photo**:
```bash
termux-camera-photo -c [0=back, 1=front] output.jpg
```
Example: `termux-camera-photo -c 0 my_photo.jpg`

### Battery
**Get status**:
```bash
termux-battery-status
```
Returns JSON with percentage, health, and plugged status.

### Torch (Flashlight)
**Turn on/off**:
```bash
termux-torch [on|off]
```

### Location
**Get current location**:
```bash
termux-location
```
Returns JSON with latitude, longitude, and accuracy.

### Microphone
**Record audio**:
```bash
termux-microphone-record -d [duration_in_seconds] -f [filename]
```
Example: `termux-microphone-record -d 5 -f recording.mp3`

## System Integration

### Clipboard
**Get content**:
```bash
termux-clipboard-get
```
**Set content**:
```bash
termux-clipboard-set "text to copy"
```

### Notifications
**Show notification**:
```bash
termux-notification --title "My Title" --content "My Message"
```

### Text-to-Speech (TTS)
**Speak text**:
```bash
termux-tts-speak [-l language] [-p pitch] [-r rate] "text to speak"
```

**Parameters**:
- `-l [language]`: Language to use (e.g., `en`, `fr`, `es`, `de`).
- `-r [rate]`: Speech rate (default: 1.0).
- `-p [pitch]`: Speech pitch (default: 1.0).

**Examples**:
Speak in English:
```bash
termux-tts-speak -l en "Hello, I am Picoclaw."
```
Speak in French:
```bash
termux-tts-speak -l fr "Bonjour, je suis Picoclaw."
```

### WiFi
**Get connection info**:
```bash
termux-wifi-connectioninfo
```
**Scan for networks**:
```bash
termux-wifi-scaninfo
```

### SMS & Telephony
**Send SMS**:
```bash
termux-sms-send -n [phone_number] "Message content"
```
**Get Call Log**:
```bash
termux-call-log -l [limit]
```

## Setup Verification

To check if the API is working:
```bash
termux-battery-status
```
If this hangs or fails, the user might need to grant permissions to the Termux:API app on their phone.
