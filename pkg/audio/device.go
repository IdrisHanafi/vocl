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

// DeviceInfoFormatter handles formatting and printing of device information
type DeviceInfoFormatter struct{}

// NewDeviceInfoFormatter creates a new device info formatter
func NewDeviceInfoFormatter() *DeviceInfoFormatter {
	return &DeviceInfoFormatter{}
}

// FormatDeviceInfo returns a formatted string with device information
func (f *DeviceInfoFormatter) FormatDeviceInfo(index int, dev *portaudio.DeviceInfo) string {
	return fmt.Sprintf("Device #%d\n  Name: %s\n  Input channels:  %d\n  Output channels: %d\n  Sample rate:     %.0f Hz\n  Input latency:   %.3f sec\n  Output latency:  %.3f sec",
		index,
		dev.Name,
		dev.MaxInputChannels,
		dev.MaxOutputChannels,
		dev.DefaultSampleRate,
		float64(dev.DefaultLowInputLatency)/float64(time.Second),
		float64(dev.DefaultLowOutputLatency)/float64(time.Second))
}

// DeviceManager handles audio device operations
type DeviceManager struct {
	devices   []*portaudio.DeviceInfo
	formatter *DeviceInfoFormatter
}

// NewDeviceManager creates a new device manager
func NewDeviceManager() (*DeviceManager, error) {
	devices, err := portaudio.Devices()
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}
	return &DeviceManager{
		devices:   devices,
		formatter: NewDeviceInfoFormatter(),
	}, nil
}

// GetDeviceByType returns a device by its ID and type (input/output)
func (dm *DeviceManager) GetDeviceByType(id int, deviceType string) (*portaudio.DeviceInfo, error) {
	// Get all devices of the specified type
	var validDevices []*portaudio.DeviceInfo
	for _, dev := range dm.devices {
		if deviceType == "input" && dev.MaxInputChannels > 0 {
			validDevices = append(validDevices, dev)
		} else if deviceType == "output" && dev.MaxOutputChannels > 0 {
			validDevices = append(validDevices, dev)
		}
	}

	if id < 0 || id >= len(validDevices) {
		return nil, fmt.Errorf("invalid device ID: %d", id)
	}

	return validDevices[id], nil
}

// GetDevice returns a device by its ID
func (dm *DeviceManager) GetDevice(id int, deviceType string) (*portaudio.DeviceInfo, error) {
	return dm.GetDeviceByType(id, deviceType)
}

// GetDefaultInputDevice returns the default input device
func (dm *DeviceManager) GetDefaultInputDevice() (*portaudio.DeviceInfo, error) {
	device, err := portaudio.DefaultInputDevice()
	if err != nil {
		return nil, err
	}
	if device.MaxInputChannels == 0 {
		return nil, fmt.Errorf("default input device '%s' has no input channels", device.Name)
	}
	return device, nil
}

// GetDefaultOutputDevice returns the default output device
func (dm *DeviceManager) GetDefaultOutputDevice() (*portaudio.DeviceInfo, error) {
	device, err := portaudio.DefaultOutputDevice()
	if err != nil {
		return nil, err
	}
	if device.MaxOutputChannels == 0 {
		return nil, fmt.Errorf("default output device '%s' has no output channels", device.Name)
	}
	return device, nil
}

// PrintDeviceInfo prints information about a device
func (dm *DeviceManager) PrintDeviceInfo(index int, dev *portaudio.DeviceInfo) {
	fmt.Println(dm.formatter.FormatDeviceInfo(index, dev))
	fmt.Println()
}

// SelectDevice interactively selects a device of the specified type
func (dm *DeviceManager) SelectDevice(deviceType string) (*portaudio.DeviceInfo, error) {
	reader := bufio.NewReader(os.Stdin)

	// Filter devices based on type
	var filteredDevices []*portaudio.DeviceInfo

	for _, dev := range dm.devices {
		if deviceType == "input" && dev.MaxInputChannels > 0 {
			filteredDevices = append(filteredDevices, dev)
		} else if deviceType == "output" && dev.MaxOutputChannels > 0 {
			filteredDevices = append(filteredDevices, dev)
		}
	}

	if len(filteredDevices) == 0 {
		return nil, fmt.Errorf("no %s devices found", deviceType)
	}

	fmt.Printf("\nAvailable %s devices:\n", deviceType)
	fmt.Println(strings.Repeat("-", len(deviceType)+20)) // Add separator line
	for i, dev := range filteredDevices {
		formatter := NewDeviceInfoFormatter()
		fmt.Println(formatter.FormatDeviceInfo(i, dev))
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
