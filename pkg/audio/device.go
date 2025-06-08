package audio

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gordonklaus/portaudio"
)

// DeviceManager handles audio device operations
type DeviceManager struct {
	devices []*portaudio.DeviceInfo
}

// NewDeviceManager creates a new device manager
func NewDeviceManager() (*DeviceManager, error) {
	devices, err := portaudio.Devices()
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}
	return &DeviceManager{devices: devices}, nil
}

// GetDevice returns a device by its ID
func (dm *DeviceManager) GetDevice(id int) (*portaudio.DeviceInfo, error) {
	if id < 0 || id >= len(dm.devices) {
		return nil, fmt.Errorf("invalid device ID: %d", id)
	}
	return dm.devices[id], nil
}

// GetDefaultInputDevice returns the default input device
func (dm *DeviceManager) GetDefaultInputDevice() (*portaudio.DeviceInfo, error) {
	return portaudio.DefaultInputDevice()
}

// GetDefaultOutputDevice returns the default output device
func (dm *DeviceManager) GetDefaultOutputDevice() (*portaudio.DeviceInfo, error) {
	return portaudio.DefaultOutputDevice()
}

// PrintDeviceInfo prints information about a device
func (dm *DeviceManager) PrintDeviceInfo(index int, dev *portaudio.DeviceInfo) {
	fmt.Printf("[%d] %s\n", index, dev.Name)
	fmt.Printf("  Max input channels:  %d\n", dev.MaxInputChannels)
	fmt.Printf("  Max output channels: %d\n", dev.MaxOutputChannels)
	fmt.Printf("  Default sample rate: %.0f Hz\n", dev.DefaultSampleRate)
	fmt.Printf("  Default low input latency:  %.3f sec\n", float64(dev.DefaultLowInputLatency)/float64(time.Second))
	fmt.Printf("  Default low output latency: %.3f sec\n", float64(dev.DefaultLowOutputLatency)/float64(time.Second))
	fmt.Println()
}

// SelectDevice interactively selects a device of the specified type
func (dm *DeviceManager) SelectDevice(deviceType string) (*portaudio.DeviceInfo, error) {
	reader := bufio.NewReader(os.Stdin)

	// Filter devices based on type
	var filteredDevices []*portaudio.DeviceInfo
	var deviceIndices []int

	for i, dev := range dm.devices {
		if deviceType == "input" && dev.MaxInputChannels > 0 {
			filteredDevices = append(filteredDevices, dev)
			deviceIndices = append(deviceIndices, i)
		} else if deviceType == "output" && dev.MaxOutputChannels > 0 {
			filteredDevices = append(filteredDevices, dev)
			deviceIndices = append(deviceIndices, i)
		}
	}

	if len(filteredDevices) == 0 {
		return nil, fmt.Errorf("no %s devices found", deviceType)
	}

	fmt.Printf("\nAvailable %s devices:\n", deviceType)
	for i, dev := range filteredDevices {
		fmt.Printf("[%d] %s\n", i, dev.Name)
		fmt.Printf("  Max %s channels: %d\n", deviceType,
			map[string]int{"input": dev.MaxInputChannels, "output": dev.MaxOutputChannels}[deviceType])
		fmt.Printf("  Default sample rate: %.0f Hz\n", dev.DefaultSampleRate)
		fmt.Printf("  Default low %s latency: %.3f sec\n", deviceType,
			map[string]float64{
				"input":  float64(dev.DefaultLowInputLatency) / float64(time.Second),
				"output": float64(dev.DefaultLowOutputLatency) / float64(time.Second),
			}[deviceType])
		fmt.Println()
	}

	for {
		fmt.Printf("Select %s device (enter number): ", deviceType)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		index, err := strconv.Atoi(input)
		if err != nil || index < 0 || index >= len(filteredDevices) {
			fmt.Println("Invalid selection. Please try again.")
			continue
		}

		return filteredDevices[index], nil
	}
}
