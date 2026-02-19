package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

const APP_VERSION = "0.0.1"

// versionCmd represents the version command
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "ascii-gen version",
		Long:  "ascii-gen version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("version %s\n", APP_VERSION)
		},
	}
}
