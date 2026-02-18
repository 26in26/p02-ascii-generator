package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render an image/video to ASCII art",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("render called")
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().StringP("output", "o", "", "Output file path")
	renderCmd.Flags().String("output-format", "terminal", "Format for output")

	renderCmd.Flags().String("resize", "100X100", "Output new dimensions")
	renderCmd.Flags().IntP("width", "w", 100, "")
	renderCmd.Flags().IntP("height", "h", 100, "")
	renderCmd.Flags().Bool("aspect-ratio", false, "")

	edge := true
	renderCmd.Flags().BoolVar(&edge, "edge", true, "Enable edge detection")
	renderCmd.Flags().BoolVar(&edge, "no-edge", true, "Disable edge detection")

	renderCmd.Flags().String("color", "none", "Select color mode")
	renderCmd.Flags().String("charset", "standard", "")
}
