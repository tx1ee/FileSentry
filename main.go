/*
Copyright © 2022 tx1ee root@tx1ee.cc

*/
package main

import (
	"filesentry/cmd"
	"filesentry/monitor"
	"github.com/gookit/color"
	"os"
	"time"
)

func main() {
	cmd.Execute()
}

func init() {
	// 判断配置文件是否存在 生成默认配置文件
	monitor.PathExists("config.yaml")
	monitor.GuserInfo, monitor.Gpushconf, monitor.Gbfilesuffix, monitor.Gmdirs, monitor.Gedirs = monitor.ReadConf()
	// 程序初始化创建目录
	initdir, _ := monitor.Pathexistence("./samples")
	if initdir {
		color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), "  samples路径初始化完成!", "\n")
	} else {
		err := os.Mkdir("./samples", os.ModePerm)
		if err != nil {
			color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), "  samples目录创建异常!", "\n")
		} else {
			color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), "  samples目录创建成功!", "\n")
		}
	}
	initdir, _ = monitor.Pathexistence("./logs")
	if initdir {
		color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), "  logs路径初始化完成!", "\n")
	} else {
		err := os.Mkdir("./logs", os.ModePerm)
		if err != nil {
			color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), "  logs目录创建异常!", "\n")
		} else {
			color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), "  logs目录创建成功!", "\n")
		}
	}
}
