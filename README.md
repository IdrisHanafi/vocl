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

#### Using Make

The project includes a Makefile for easy building and installation:

```bash
# Build the application
make build

# Install to /usr/local/bin
make install

# Remove build artifacts
make clean

# Uninstall from /usr/local/bin
make uninstall

# Show available commands
make help
```

#### Manual Build

Alternatively, you can build manually:

```bash
go build -o build/vocl main.go
```

This will create an executable named `vocl` in the `build` directory.

## Usage

### Available Commands

The tool provides two main commands:

1. `info` - List all available audio devices
2. `run` - Start real-time audio processing with customizable effects

### Listing Audio Devices

To see all available input and output devices:

```bash
./vocl info
```

This will show a list of all audio devices with their IDs, which you can use with the `-i` and `-o` flags. The output looks like this:

```
Input Devices:
-------------
Device #0
  Name: MacBook Pro Microphone
  Input channels:  1
  Output channels: 0
  Sample rate:     48000 Hz
  Input latency:   0.032 sec
  Output latency:  0.010 sec

Output Devices:
--------------
Device #0
  Name: MacBook Pro Speakers
  Input channels:  0
  Output channels: 2
  Sample rate:     48000 Hz
  Input latency:   0.010 sec
  Output latency:  0.019 sec
```

Use the device numbers (e.g., `Device #0`) as the values for the `-i` and `-o` flags.

### Running Audio Processing

Run the tool with default settings:

```bash
./vocl run
```

### Command-line Options

The `run` command supports the following flags:

```bash
./vocl run [flags]

Flags:
  -d, --delay int       Echo delay time in milliseconds (how long before the echo repeats) (default 300)
  -f, --feedback float  Echo feedback (0-1) (how much of the echo is fed back into itself, creating multiple repeats) (default 0.7)
  -m, --mix float       Echo mix level (0-1) (how loud the echo is compared to the original sound) (default 0.7)
  -i, --input int       Input device ID (use -1 for default input device)
  -o, --output int      Output device ID (use -1 for default output device)
```

### Device Selection

The tool provides flexible device selection:

1. Run without device flags to select devices interactively:
   ```bash
   ./vocl run
   ```

2. Specify input device only (interactive output selection):
   ```bash
   ./vocl run -i 0
   ```

3. Specify output device only (interactive input selection):
   ```bash
   ./vocl run -o 1
   ```

4. Specify both devices:
   ```bash
   ./vocl run -i 0 -o 1
   ```

5. Use default devices:
   ```bash
   ./vocl run -i -1 -o -1
   ```

When no device is specified, the tool will:
1. Show a list of available input devices
2. Prompt you to select an input device
3. Show a list of available output devices
4. Prompt you to select an output device
5. Start streaming audio with the selected effects

# Troubleshoot

If you encounter issues with audio input, make sure you have the necessary permissions to access the microphone and speaker on your system. You can check this in your system settings.

Mac users may need to allow terminal access to the microphone in System Preferences > Security & Privacy > Privacy > Microphone.

# Future Work

Feel free to create a PR on any of the following tasks:

- [ ] Filter background noise (like puffs, statics, etc.). Can't be bothered to understand the math on this right now
