# Overview

This is a lightweight CLI tool written in golang to setup a karaoke speaker
using your system's microphone and speaker.

## Setup

### Prerequisites

- Go 1.20 or later
- [PortAudio](https://www.portaudio.com/)
- [pkg-config][https://www.freedesktop.org/wiki/Software/pkg-config/](Tells `Go` where to find `C` libraries)

### Build Locally

On macOS, you can install `PortAudio` and `pkg-config` using Homebrew:

```bash
brew install pkg-config portaudio
```

Minimally, you can run:
```bash
go run main.go
```

# Troubleshoot

If you encounter issues with audio input, make sure you have the necessary permissions to access the microphone and speaker on your system. You can check this in your system settings.

Mac users may need to allow terminal access to the microphone in System Preferences > Security & Privacy > Privacy > Microphone.
