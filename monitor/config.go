package monitor

import (
	"github.com/gookit/color"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// 用户信息服务器信息
type UserInfo struct {
	User        string   `yaml:"User"`
	ServerName  string   `yaml:"ServerName"`
	ServerIP    string   `yaml:"ServerIP"`
	ServerMac   string   `yaml:"ServerMac"`
	PushConf    PushConf `yaml:"PushConf"`
	Bfilesuffix []string `yaml:"Bfilesuffix"`
	MonitorDirs []string `yaml:"MonitorDirs"`
	ExcludeDirs []string `yaml:"ExcludeDirs"`
}

// 消息推送配置
type PushConf struct {
	DingToken  string `yaml:"DingToken"`
	DingSecert string `yaml:"DingSecert"`
	MailUser   string `yaml:"MailUser"`
	MailKey    string `yaml:"MailKey"`
	Mailadder  string `yaml:"Mailaddr"`
	// 安全检测
	SecCheck bool `yaml:"SecCheck"`
	// 邮件推送
	Mailpush bool `yaml:"Mailpush"`
}

// 写入yaml
func writeToXml(src string, UserInfos []string, Bfilesuffix []string, PushConfs PushConf, MDirs []string, EDir []string) {
	stu := &UserInfo{
		User:        UserInfos[0],
		ServerName:  UserInfos[1],
		ServerIP:    UserInfos[2],
		ServerMac:   UserInfos[3],
		Bfilesuffix: Bfilesuffix,
		PushConf:    PushConfs,
		MonitorDirs: MDirs,
		ExcludeDirs: EDir,
	}
	data, err := yaml.Marshal(stu)
	checkError(err)
	err = ioutil.WriteFile(src, data, 0777)
	checkError(err)
}

func ReadConf() ([4]string, [7]string, []string, []string, []string) {
	//获取项目的执行路径
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	config := viper.New()
	config.AddConfigPath(path)     //设置读取的文件路径
	config.SetConfigName("config") //设置读取的文件名
	config.SetConfigType("yaml")   //设置文件的类型
	//尝试进行配置读取
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	var userinfo [4]string
	var pushconf [7]string
	// 用户信息配置
	userinfo[0] = config.GetString("User")
	//fmt.Println(userinfo[0])
	userinfo[1] = config.GetString("ServerName")
	userinfo[2] = config.GetString("ServerIP")
	userinfo[3] = config.GetString("ServerMac")
	// 消息推送配置
	pushconf[0] = config.GetString("PushConf.DingToken")
	pushconf[1] = config.GetString("PushConf.DingSecert")
	pushconf[2] = config.GetString("PushConf.MailUser")
	pushconf[3] = config.GetString("PushConf.MailKey")
	pushconf[4] = config.GetString("PushConf.Mailaddr")
	pushconf[5] = config.GetString("PushConf.SecCheck")
	pushconf[6] = config.GetString("PushConf.Mailpush")

	// 黑名单后缀
	bfilesuffix := config.GetStringSlice("Bfilesuffix")
	// 目录监控配置
	mdirs := config.GetStringSlice("MonitorDirs")
	edirs := config.GetStringSlice("ExcludeDirs")
	return userinfo, pushconf, bfilesuffix, mdirs, edirs
}

// yaml配置文件修改
func WritePushConf(pushconf []string) {
	var data []byte
	file := "config.yaml"
	// 读取文件
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Print(err)
		return
	}
	var conf UserInfo
	// 转换成Struct
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		Warning.Println("%v\n", err.Error())
	}
	// 修改pushconf
	conf.PushConf.DingToken = pushconf[0]
	conf.PushConf.DingSecert = pushconf[1]
	conf.PushConf.MailUser = pushconf[2]
	conf.PushConf.MailKey = pushconf[3]
	conf.PushConf.Mailadder = pushconf[4]
	data, err = yaml.Marshal(&conf)
	if err != nil {
		Warning.Println("%v\n", err.Error())
	}
	err = ioutil.WriteFile("config.yaml", data, 0777)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		InitConf()
		return false, nil
	}
	return false, err
}

func InitConf() {
	ip := GetIps()
	var MonitorDir []string
	var ExcludeDir []string
	var UserInfos = []string{"测试客户", "测试服务", ip[0], "00-00-00-00-00-00"}
	if runtime.GOOS == "linux" {
		color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), " 当前服务器系统为linux\n")
		Info.Println("当前服务器系统为linux")
		MonitorDir = append(MonitorDir, "/home/tx1ee")
		ExcludeDir = append(ExcludeDir, "/root/tx1ee/test")
	} else if runtime.GOOS == "windows" {
		Info.Println("当前服务器系统为windows")
		MonitorDir = append(MonitorDir, "c:\\")
		ExcludeDir = append(ExcludeDir, "c:\\tx1ee\\test")
	}
	var PushConfs = PushConf{"you dingding token", "you dingding secret", "you mail account", "you mail password", "Email notification recipient", true, true}
	var Bfilesuffix = []string{".php", ".exe", ".php3"}
	writeToXml("config.yaml", UserInfos, Bfilesuffix, PushConfs, MonitorDir, ExcludeDir)
}
