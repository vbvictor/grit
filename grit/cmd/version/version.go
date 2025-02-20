package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "v0.1.0"

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Args:  cobra.ExactArgs(0),
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Printf("Version: %s", Version)

		return nil
	},
}
