package vocl

import (
	"fmt"
	"log"

	"github.com/IdrisHanafi/vocl/pkg/audio"
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

		dm, err := audio.NewDeviceManager()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("\nInput Devices:")
		fmt.Println("-------------")
		inputIndex := 0
		for _, dev := range devices {
			if dev.MaxInputChannels > 0 {
				dm.PrintDeviceInfo(inputIndex, dev)
				inputIndex++
			}
		}

		fmt.Println("\nOutput Devices:")
		fmt.Println("--------------")
		outputIndex := 0
		for _, dev := range devices {
			if dev.MaxOutputChannels > 0 {
				dm.PrintDeviceInfo(outputIndex, dev)
				outputIndex++
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
