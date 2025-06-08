package vocl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vocl",
	Short: "VOCL - Voice Over Command Line",
	Long: `VOCL is a command-line tool for audio processing and effects.
It allows you to process audio input with various effects like echo,
and provides device information for your audio setup.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
