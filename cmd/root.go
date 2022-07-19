/*
Copyright © 2022 tx1ee root@tx1ee.cc
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golangcode/monitor"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "sentry",
	CompletionOptions: struct {
		DisableDefaultCmd   bool
		DisableNoDescFlag   bool
		DisableDescriptions bool
		HiddenDefaultCmd    bool
	}{DisableDefaultCmd: true, DisableNoDescFlag: false, DisableDescriptions: false, HiddenDefaultCmd: false},
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "run" {
			for _, gmdir := range monitor.Gmdirs {
				// 判断要监控的路径是否存在 不存在程序会报错 路径不存在直接结束程序
				mdircz, _ := monitor.Pathexistence(gmdir)
				if mdircz {
					watch := monitor.NewNotifyFile()
					watch.WatchDir(gmdir, monitor.Gedirs, monitor.Gbfilesuffix)
				} else {
					monitor.Warning.Println("你监控的" + gmdir + "路径不存在，程序退出...")
					os.Exit(0)
				}

			}
			select {}
			return
		} else if args[0] == "push" {
			// 测试推送消息是否正常
			monitor.Dingpushtest(monitor.Gpushconf[0], monitor.Gpushconf[1])
			monitor.Info.Println("消息推送测试")
			//fmt.Println("消息推送测试")
		}
	},
}

func init() {
	fmt.Println(monitor.Logo1)
}

func Execute() {
	rootCmd.Execute()
}
