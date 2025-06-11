package audio

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
)

// StreamConfig holds the configuration for an audio stream
type StreamConfig struct {
	InputDevice  *portaudio.DeviceInfo
	OutputDevice *portaudio.DeviceInfo
	SampleRate   float64
	Channels     int
	BufferSize   int
}

// Stream represents an audio stream with effects
type Stream struct {
	stream  *portaudio.Stream
	config  StreamConfig
	effects []Effect
}

// NewStream creates a new audio stream with the given configuration and effects
func NewStream(config StreamConfig, effects []Effect) (*Stream, error) {
	if len(effects) != config.Channels {
		return nil, fmt.Errorf("number of effects (%d) must match number of channels (%d)",
			len(effects), config.Channels)
	}

	stream, err := portaudio.OpenStream(portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   config.InputDevice,
			Channels: config.Channels,
			Latency:  config.InputDevice.DefaultLowInputLatency,
		},
		Output: portaudio.StreamDeviceParameters{
			Device:   config.OutputDevice,
			Channels: config.Channels,
			Latency:  config.OutputDevice.DefaultLowOutputLatency,
		},
		SampleRate:      config.SampleRate,
		FramesPerBuffer: config.BufferSize,
	}, func(inBuf, outBuf []float32, _ portaudio.StreamCallbackTimeInfo, _ portaudio.StreamCallbackFlags) {
		// Process each channel through its effect
		// TODO: Future optimization to maybe parallelize this processing for each channel
		for i := 0; i < config.Channels; i++ {
			// Get the samples for this channel
			channelIn := make([]float32, len(inBuf)/config.Channels)
			channelOut := make([]float32, len(outBuf)/config.Channels)

			// Extract channel data
			for j := 0; j < len(channelIn); j++ {
				channelIn[j] = inBuf[j*config.Channels+i]
			}

			// Process through effect
			for j := 0; j < len(channelIn); j++ {
				channelOut[j] = effects[i].Process(channelIn[j])
			}

			// Put processed data back
			for j := 0; j < len(channelOut); j++ {
				outBuf[j*config.Channels+i] = channelOut[j]
			}
		}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to open stream: %w", err)
	}

	return &Stream{
		stream:  stream,
		config:  config,
		effects: effects,
	}, nil
}

// Start starts the audio stream
func (s *Stream) Start() error {
	return s.stream.Start()
}

// Stop stops the audio stream
func (s *Stream) Stop() error {
	return s.stream.Stop()
}

// Close closes the audio stream
func (s *Stream) Close() error {
	return s.stream.Close()
}
