package vocl

import (
	"fmt"
	"log"

	"github.com/IdrisHanafi/vocl/pkg/audio"
	"github.com/gordonklaus/portaudio"
	"github.com/spf13/cobra"
)

type EchoConfig struct {
	DelayTime float32
	Feedback  float32
	Mix       float32
}

var (
	inputDeviceID  int
	outputDeviceID int
	delayTime      float32
	feedback       float32
	mix            float32
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run audio processing with echo effect",
	Long:  `Start audio processing with configurable echo effect parameters. Use flags to set parameters or interactive mode for device selection.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize PortAudio
		if err := portaudio.Initialize(); err != nil {
			log.Fatal(err)
		}
		defer portaudio.Terminate()

		// Create device manager
		dm, err := audio.NewDeviceManager()
		if err != nil {
			log.Fatal(err)
		}

		var inputDevice, outputDevice *portaudio.DeviceInfo

		// Handle input device selection
		if cmd.Flags().Changed("input") {
			// Use specified input device or default if -1
			if inputDeviceID < 0 {
				inputDevice, err = dm.GetDefaultInputDevice()
			} else {
				inputDevice, err = dm.GetDevice(inputDeviceID, "input")
			}
		} else {
			// No input device specified, use interactive selection
			inputDevice, err = dm.SelectDevice("input")
		}
		if err != nil {
			log.Fatal(err)
		}

		// Handle output device selection
		if cmd.Flags().Changed("output") {
			// Use specified output device or default if -1
			if outputDeviceID < 0 {
				outputDevice, err = dm.GetDefaultOutputDevice()
			} else {
				outputDevice, err = dm.GetDevice(outputDeviceID, "output")
			}
		} else {
			// No output device specified, use interactive selection
			outputDevice, err = dm.SelectDevice("output")
		}
		if err != nil {
			log.Fatal(err)
		}

		// Use the minimum number of channels between input and output
		channels := inputDevice.MaxInputChannels
		channels = min(outputDevice.MaxOutputChannels, channels)

		fmt.Printf("\nSelected input device: %s\n", inputDevice.Name)
		fmt.Printf("Selected output device: %s\n", outputDevice.Name)
		fmt.Printf("Using sample rate: %.0f Hz\n", inputDevice.DefaultSampleRate)
		fmt.Printf("Using channels: %d\n", channels)
		fmt.Printf("Echo settings - Delay: %.0fms, Feedback: %.2f, Mix: %.2f\n",
			delayTime, feedback, mix)

		// Create echo effects for each channel
		echoEffects := make([]audio.Effect, inputDevice.MaxInputChannels)
		for i := range echoEffects {
			echoEffects[i] = audio.NewEchoEffect(float64(inputDevice.DefaultSampleRate), delayTime, feedback, mix)
		}

		// Create and start the stream
		stream, err := audio.NewStream(audio.StreamConfig{
			InputDevice:  inputDevice,
			OutputDevice: outputDevice,
			SampleRate:   inputDevice.DefaultSampleRate,
			Channels:     channels,
			BufferSize:   64,
		}, echoEffects)
		if err != nil {
			log.Fatal(err)
		}
		defer stream.Close()

		if err := stream.Start(); err != nil {
			log.Fatal(err)
		}

		log.Println("Streaming mic input to output with echo effect. Press Ctrl+C to stop.")
		select {} // block forever
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Add flags
	runCmd.Flags().IntVarP(&inputDeviceID, "input", "i", -1, "Input device ID (use -1 for default input device)")
	runCmd.Flags().IntVarP(&outputDeviceID, "output", "o", -1, "Output device ID (use -1 for default output device)")
	runCmd.Flags().Float32VarP(&delayTime, "delay", "d", 300, "Echo delay time in milliseconds (how long before the echo repeats)")
	runCmd.Flags().Float32VarP(&feedback, "feedback", "f", 0.7, "Echo feedback (0-1) (how much of the echo is fed back into itself, creating multiple repeats)")
	runCmd.Flags().Float32VarP(&mix, "mix", "m", 0.7, "Echo mix level (0-1) (how loud the echo is compared to the original sound)")
}
