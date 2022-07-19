/*
Copyright © 2022 tx1eee root@tx1ee.cc

*/
package cmd

import (
	"github.com/spf13/cobra"
	"golangcode/monitor"
)

var PushConf string

var confCmd = &cobra.Command{
	Use:   "conf",
	Short: "Configure message push",
	Args:  cobra.MinimumNArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		monitor.WritePushConf(args)
	},
}

func init() {
	rootCmd.AddCommand(confCmd)
	rootCmd.SetHelpTemplate(
		"Usage:\n" +
			"  sentry conf dingtoken dingsecert mailuser mailkey mailaddr\n" +
			"  sentry conf 82*****895ce SEC***2 root@tx1ee.cc ED****LK admin@tx1ee.cc\n\n" +
			"Flags:\n  -h, --help   help for conf\n")
	rootCmd.SetUsageTemplate("Usage:\n" +
		"" +
		"  sentry conf dingtoken dingsecert mailuser mailkey mailaddr\n" +
		"  sentry conf 82*****895ce SEC***2 root@tx1ee.cc ED****LK admin@tx1ee.cc\n\n")
}
