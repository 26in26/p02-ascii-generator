package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display application information",
	Long:  `Display detailed information about the ASCII Generator application, including features and version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nASCII Generator (ascii-gen)")
		fmt.Printf("Version: %s\n\n", APP_VERSION)
		fmt.Println("Description:")
		fmt.Println("  A high-performance CLI tool designed to convert images into ASCII art.")
		fmt.Println("  It utilizes a multi-stage pipeline including resizing, grayscale conversion,")
		fmt.Println("  Sobel edge detection, and character mapping to produce high-quality text representations.")
		fmt.Println("\nKey Features:")
		fmt.Println("  - Optimized Performance: Uses integer arithmetic and bitwise operations.")
		fmt.Println("  - Edge Detection: Uses Sobel operator to capture image structure.")
		fmt.Println("  - Color Modes: Supports full 24-bit TrueColor and monochrome output.")
		fmt.Println("  - Configurable: Adjustable dimensions, thresholds, and character sets.")
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
