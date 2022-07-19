package monitor

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//全局变量
var (
	GuserInfo    [4]string
	Gpushconf    [7]string
	Gbfilesuffix []string
	Gmdirs       []string
	Gedirs       []string
	Dingpushinfo [9]string
)

// 消息推送

// Pathexistence 判断所给路径文件/文件夹是否存在
func Pathexistence(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//获取本机ip地址
func GetIps() (ips []string) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interfaces ipAddress: %v\n", err)
		return ips
	}

	for _, address := range interfaceAddr {
		ipNet, isVailIpNet := address.(*net.IPNet)
		if isVailIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

//Logo
var Logo = string("\n" +
	"               /―――――――――――――/" + "\n" +
	"              /  Blue team  /" + "       (๑°⌓°๑)" + "\n" +
	"    _?_?_    /―――――――――――――/" + "\n" +
	"    &) &)   /              _                 \n" +
	"  ---/⌓/---/  ____  ____  | |_    ____  _   _ \n" +
	"     \\ \\     / _  )|  _ \\ |  _)  / ___)| | | |\n" +
	" _____) )   ( (/ / | | | || |__ | |    | |_| |\n" +
	"(______/     \\____)|_| |_| \\___)|_|     \\__  |\n" +
	"                                       (_____/  v0.1 by tx1ee   \n" +
	"\n" +
	"工具介绍: 带消息推送和安全检测的文件监控工具,辅助蓝队监控web网站目录情况\n" +
	"github: https://github.com/tx1ee\n")

// Logo1
var Logo1 = string("\n" +
	"                                     /-------------------/" + "\n" +
	"                                    /     Blue team     /" + "\n" +
	" _______  _  _             _?_?_   /-------------------/" + "\n" +
	"(_______)(_)| |          -*) -*)  /             _                 \n" +
	" _____    _ | |  ____   ---/o/---/ ____  ____  | |_    ____  _   _ \n" +
	"|  ___)  | || | / _  )     \\ \\    / _  )|  _ \\ |  _)  / ___)| | | |\n" +
	"| |      | || |( (/ /  _____) )  ( (/ / | | | || |__ | |    | |_| |\n" +
	"|_|      |_||_| \\____)(______/    \\____)|_| |_| \\___)|_|     \\__  |\n" +
	"                                                             (_____/  v0.1 by tx1ee   \n" +
	"\n" +
	"工具介绍: 带消息推送和安全检测的文件监控工具,辅助蓝队监控web网站目录情况\n" +
	"github: https://github.com/tx1ee\n")

// RandString 随机字符生成事件ID
func RandString(length int) string {
	rand.Seed(time.Now().UnixNano())
	rs := make([]string, length)
	for start := 0; start < length; start++ {
		t := rand.Intn(3)
		if t == 0 {
			rs = append(rs, strconv.Itoa(rand.Intn(10)))
		} else if t == 1 {
			rs = append(rs, string(rand.Intn(26)+65))
		} else {
			rs = append(rs, string(rand.Intn(26)+97))
		}
	}
	return strings.Join(rs, "")
}
