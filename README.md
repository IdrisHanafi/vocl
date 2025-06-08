# Overview

This is a lightweight CLI tool written in golang to setup a karaoke speaker
using your system's microphone and speaker.

## Setup

### Prerequisites

- Go 1.20 or later
- [PortAudio](https://www.portaudio.com/)
- [pkg-config](https://www.freedesktop.org/wiki/Software/pkg-config/) (Tells `Go` where to find `C` libraries)

### Build Locally

On macOS, you can install `PortAudio` and `pkg-config` using Homebrew:

```bash
brew install pkg-config portaudio
```

To build the CLI tool:

```bash
go build -o vocl cmd/vocl/main.go
```

This will create an executable named `vocl` in your current directory.

## Usage

Run the tool with default settings:

```bash
./vocl run
```

### Command-line Options

The `run` command supports the following flags:

```bash
./vocl run [flags]

Flags:
  -d, --delay int       Echo delay time in milliseconds (default 100)
  -f, --feedback float  Echo feedback amount (0.0 to 1.0) (default 0.2)
  -m, --mix float       Echo mix level (0.0 to 1.0) (default 0.6)
  -i, --input int       Input device ID
  -o, --output int      Output device ID
```

### Interactive Mode

If you don't specify input/output devices, the tool will enter interactive mode and prompt you to select your devices:

```bash
./vocl run
```

This will:
1. Show a list of available input devices
2. Prompt you to select an input device
3. Show a list of available output devices
4. Prompt you to select an output device
5. Start streaming audio with echo effect

# Troubleshoot

If you encounter issues with audio input, make sure you have the necessary permissions to access the microphone and speaker on your system. You can check this in your system settings.

Mac users may need to allow terminal access to the microphone in System Preferences > Security & Privacy > Privacy > Microphone.

# Future Work

Feel free to create a PR on any of the following tasks:

- [ ] Filter background noise (like puffs, statics, etc.). Can't be bothered to understand the math on this right now