package cmd

import (
	"os"

	"github.com/26in26/p02-ascii-generator/internal/cli/info"
	"github.com/26in26/p02-ascii-generator/internal/cli/render"
	"github.com/26in26/p02-ascii-generator/internal/cli/version"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ascii-gen",
	Short: "A high-performance CLI tool for converting images to ASCII art",
	Long: `ASCII Generator is a blazing fast, highly configurable command-line tool that transforms your images into stunning ASCII art.

Built with performance in mind, it utilizes advanced image processing techniques including:
  - Sobel Edge Detection for crisp structural details
  - 24-bit TrueColor support for vibrant output
  - Optimized integer arithmetic for maximum speed

Whether you want retro monochrome aesthetics or full-color text representations, ascii-gen delivers high-quality results instantly.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(info.NewCommand())
	rootCmd.AddCommand(render.NewCommand())
	rootCmd.AddCommand(version.NewCommand())

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
