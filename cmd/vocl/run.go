package vocl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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
	interactive    bool
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run audio processing with echo effect",
	Long:  `Start audio processing with configurable echo effect parameters. Use flags to set parameters or --interactive for interactive mode.`,
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

		if interactive {
			// Interactive mode
			inputDevice, err = dm.SelectDevice("input")
			if err != nil {
				log.Fatal(err)
			}
			outputDevice, err = dm.SelectDevice("output")
			if err != nil {
				log.Fatal(err)
			}
			config := getEchoConfig()
			delayTime = config.DelayTime
			feedback = config.Feedback
			mix = config.Mix
		} else {
			// Use flags
			if inputDeviceID < 0 {
				// Use default input device
				inputDevice, err = dm.GetDefaultInputDevice()
				if err != nil {
					log.Fatal("Failed to get default input device:", err)
				}
			} else {
				inputDevice, err = dm.GetDevice(inputDeviceID)
				if err != nil {
					log.Fatal(err)
				}
			}

			if outputDeviceID < 0 {
				// Use default output device
				outputDevice, err = dm.GetDefaultOutputDevice()
				if err != nil {
					log.Fatal("Failed to get default output device:", err)
				}
			} else {
				outputDevice, err = dm.GetDevice(outputDeviceID)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		// Use the minimum number of channels between input and output
		channels := inputDevice.MaxInputChannels
		if outputDevice.MaxOutputChannels < channels {
			channels = outputDevice.MaxOutputChannels
		}

		fmt.Printf("\nSelected input device: %s\n", inputDevice.Name)
		fmt.Printf("Selected output device: %s\n", outputDevice.Name)
		fmt.Printf("Using sample rate: %.0f Hz\n", inputDevice.DefaultSampleRate)
		fmt.Printf("Using channels: %d\n", channels)
		fmt.Printf("Echo settings - Delay: %.0fms, Feedback: %.2f, Mix: %.2f\n",
			delayTime, feedback, mix)

		// Create echo effects for each channel
		echoEffects := make([]audio.Effect, inputDevice.MaxInputChannels)
		for i := range echoEffects {
			echoEffects[i] = audio.NewEchoEffect(float64(inputDevice.DefaultSampleRate), 100, 0.2, 0.6)
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

func getEchoConfig() EchoConfig {
	reader := bufio.NewReader(os.Stdin)
	config := EchoConfig{
		DelayTime: 100,
		Feedback:  0.2,
		Mix:       0.6,
	}

	fmt.Println("\nEcho Configuration")
	fmt.Println("-----------------")

	// Get delay time
	for {
		fmt.Printf("Enter delay time in milliseconds (default: %.0f): ", config.DelayTime)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			break
		}
		if val, err := strconv.ParseFloat(input, 32); err == nil && val > 0 {
			config.DelayTime = float32(val)
			break
		}
		fmt.Println("Invalid input. Please enter a positive number.")
	}

	// Get feedback
	for {
		fmt.Printf("Enter feedback (0-1, default: %.2f): ", config.Feedback)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			break
		}
		if val, err := strconv.ParseFloat(input, 32); err == nil && val >= 0 && val <= 1 {
			config.Feedback = float32(val)
			break
		}
		fmt.Println("Invalid input. Please enter a number between 0 and 1.")
	}

	// Get mix
	for {
		fmt.Printf("Enter mix (0-1, default: %.2f): ", config.Mix)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			break
		}
		if val, err := strconv.ParseFloat(input, 32); err == nil && val >= 0 && val <= 1 {
			config.Mix = float32(val)
			break
		}
		fmt.Println("Invalid input. Please enter a number between 0 and 1.")
	}

	return config
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Add flags
	runCmd.Flags().IntVarP(&inputDeviceID, "input", "i", -1, "Input device ID (use -1 for default input device)")
	runCmd.Flags().IntVarP(&outputDeviceID, "output", "o", -1, "Output device ID (use -1 for default output device)")
	runCmd.Flags().Float32VarP(&delayTime, "delay", "d", 100, "Echo delay time in milliseconds")
	runCmd.Flags().Float32VarP(&feedback, "feedback", "f", 0.2, "Echo feedback (0-1)")
	runCmd.Flags().Float32VarP(&mix, "mix", "m", 0.6, "Echo mix level (0-1)")
	runCmd.Flags().BoolVarP(&interactive, "interactive", "t", true, "Enable interactive mode for device and parameter selection")
}
