package vocl

import (
	"fmt"
	"log"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display audio device information",
	Long:  `Display detailed information about available input and output audio devices.`,
	Run: func(cmd *cobra.Command, args []string) {
		portaudio.Initialize()
		defer portaudio.Terminate()

		devices, err := portaudio.Devices()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("\nInput Devices:")
		fmt.Println("-------------")
		for i, dev := range devices {
			if dev.MaxInputChannels > 0 {
				printDeviceInfo(i, dev)
			}
		}

		fmt.Println("\nOutput Devices:")
		fmt.Println("--------------")
		for i, dev := range devices {
			if dev.MaxOutputChannels > 0 {
				printDeviceInfo(i, dev)
			}
		}
	},
}

func printDeviceInfo(index int, dev *portaudio.DeviceInfo) {
	fmt.Printf("[%d] %s\n", index, dev.Name)
	fmt.Printf("  Max input channels:  %d\n", dev.MaxInputChannels)
	fmt.Printf("  Max output channels: %d\n", dev.MaxOutputChannels)
	fmt.Printf("  Default sample rate: %.0f Hz\n", dev.DefaultSampleRate)
	fmt.Printf("  Default low input latency:  %.3f sec\n", float64(dev.DefaultLowInputLatency)/float64(time.Second))
	fmt.Printf("  Default low output latency: %.3f sec\n", float64(dev.DefaultLowOutputLatency)/float64(time.Second))
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
