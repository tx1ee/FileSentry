package monitor

import (
	"fmt"
	"github.com/blinkbean/dingtalk"
	"github.com/gookit/color"
	"gopkg.in/gomail.v2"
	"time"
)

var (
	FromMail string
	ToMail   string
	MailUser string
	MailPwd  string
)

// SendMail 邮件样本推送
func SendMail() {
	m := gomail.NewMessage()
	//发送人
	m.SetHeader("From", FromMail)
	//接收人
	m.SetHeader("To", ToMail)
	//抄送人
	//m.SetAddressHeader("Cc", "xxx@qq.com", "xiaozhujiao")
	//主题
	m.SetHeader("Subject", "Sentry检测到恶意病毒木马")
	//内容
	m.SetBody("text/html", "<h1>病毒样本</h1>")
	//附件
	m.Attach("./test.jpg")
	//拿到token，并进行连接,第4个参数是填授权码
	d := gomail.NewDialer("smtp.qq.com", 587, MailUser, MailPwd)
	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("DialAndSend err %v:", err)
		panic(err)
	}
	// 修改
	Info.Printf("send mail success\n")
}

// Dingpush 钉钉消息推送
func Dingpush(dingtoken string, dingsecert string, dingpushinfo [9]string) {
	var dingToken = dingtoken
	cli := dingtalk.InitDingTalkWithSecret(dingToken, dingsecert)
	nowtime := time.Now().Format("2006-01-02 15:04:05")
	// 安全检测消息模板
	msg := []string{
		"# Sentry文件监控提醒",
		"---",
		"## " + nowtime + "\n",
		"---",
		"🕵️有老6?",
		"### [+] 监控信息\n",
		"客户名称：" + GuserInfo[0] + "\n",
		"服务器名：" + GuserInfo[1] + "\n",
		"服务器IP：" + GuserInfo[2] + "\n",
		"服务器MAC：" + GuserInfo[3] + "\n",
		"### [+] 事件详细信息\n",
		"事件ID：<font color=#ff2a00 size=20>" + dingpushinfo[0] + "</font>\n",
		"触发事件：<font color=#ff2a00 size=20>" + dingpushinfo[1] + "</font>\n",
		"事件时间：" + nowtime + "\n",
		"相关文件：" + dingpushinfo[2] + "\n",
		"### [+] 安全检测结果\n",
		"百度Webdir+检测结果：<font color=#ff2a00 size=20>" + dingpushinfo[3] + "</font>\n",
		"河马webshell检测结果：<font color=#ff2a00 size=20>" + dingpushinfo[4] + "</font>\n",
	}
	// 不做安全检测消息模板
	msg1 := []string{
		"# Sentry文件监控提醒",
		"---",
		"## " + nowtime + "\n",
		"---",
		"🕵️有老6?",
		"### [+] 监控信息\n",
		"客户名称：" + GuserInfo[0] + "\n",
		"服务器名：" + GuserInfo[1] + "\n",
		"服务器IP：" + GuserInfo[2] + "\n",
		"服务器MAC：" + GuserInfo[3] + "\n",
		"### [+] 事件详细信息\n",
		"事件ID：<font color=#ff2a00 size=20>" + dingpushinfo[0] + "</font>\n,",
		"触发事件：<font color=#ff2a00 size=20>" + dingpushinfo[1] + "</font>\n",
		"事件时间：" + nowtime + "\n",
		"相关文件：" + dingpushinfo[2] + "\n",
		"### [+] 安全检测结果\n",
		"[INFO] <font color=#ff2a00 size=20>" + "未开启安全检测" + "</font>\n,"}
	if Gpushconf[5] == "true" {
		// 安全检测结果模板
		cli.SendMarkDownMessageBySlice("Sentry监控事件提醒", msg, dingtalk.WithAtAll())
	} else {
		// 无安全检测结果模板
		cli.SendMarkDownMessageBySlice("Sentry监控事件提醒", msg1, dingtalk.WithAtAll())
	}
}

func Dingpushtest(dingtoken string, dingsecert string) {
	var dingToken = dingtoken
	cli := dingtalk.InitDingTalkWithSecret(dingToken, dingsecert)
	nowtime := time.Now().Format("2006-01-02 15:04:05")
	msgtest := []string{
		"# Sentry消息推送测试",
		"---",
		"## " + nowtime + "\n",
		"---",
	}
	cli.SendMarkDownMessageBySlice("Sentry监控事件提醒", msgtest, dingtalk.WithAtAll())
	color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), "  钉钉消息推送测试!", "\n")

}
