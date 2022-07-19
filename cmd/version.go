/*
Copyright © 2022 tx1ee root@tx1ee.cc

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "About version information",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("2022-07-17 Sentry v0.1 by tx1ee")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
