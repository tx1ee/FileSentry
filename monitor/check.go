package monitor

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func Hmscan(scanfile string) []string {
	// 上传到河马webshell在线检测
	url := "https://n.shellpub.com/api/v1/zip"
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()
		part, err := m.CreateFormFile("userfile", "scanfile.zip")
		if err != nil {
			return
		}
		file, err := os.Open(scanfile)
		if err != nil {
			return
		}
		defer file.Close()
		if _, err = io.Copy(part, file); err != nil {
			return
		}
	}()
	a, _ := http.Post(url, m.FormDataContentType(), r)
	rspBody, _ := ioutil.ReadAll(a.Body)
	result := viper.New()
	// 设置配置文件类型为json 这里很重要 不是指的话会读取不到 踩坑踩坑
	result.SetConfigType("json")
	result.ReadConfig(bytes.NewBuffer(rspBody))
	if err := result.ReadConfig(bytes.NewBuffer(rspBody)); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("找不到配置文件..")
		} else {
			fmt.Println("配置文件出错..")
		}
	}
	// 结果url拼接
	result_url := "https://n.shellpub.com/detail/" + result.GetString(`data.tid`)
	scanresult := hmresult(result_url)
	return scanresult
}
func hmresult(url string) []string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	// 请求头设置 设置长连接
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.0.0 Safari/537.36")
	req.Header.Set("Connection", "keep-alive")
resetcheck:
	// 开始请求
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Http get err:", err)
	}
	//if resp.StatusCode != 200 {
	//	fmt.Println("status code:", err)
	//}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if checkstatus := strings.Contains(string(body), "服务目前不可用"); checkstatus {
		time.Sleep(3 * time.Second)
		goto resetcheck
	} else {
		fmt.Printf("河马webshell扫描完成")
	}
	//fmt.Printf(string(body))
	html := strings.Replace(string(body), "\n", "", -1)
	re_json := regexp.MustCompile(`"results":\[(.*?)\]`)
	result_json := re_json.FindString(html)
	result_json = result_json[11 : len(result_json)-1]
	result := viper.New()
	// 设置配置文件类型为json 这里很重要 不是指的话会读取不到 踩坑踩坑
	result.SetConfigType("json")
	result.ReadConfig(bytes.NewBuffer([]byte(result_json)))
	if err := result.ReadConfig(bytes.NewBuffer([]byte(result_json))); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("找不到配置文件..")
		} else {
			fmt.Println("配置文件出错..")
		}
	}

	// 返回字符数组
	var scanresult []string
	scanresult = append(scanresult, result.GetString(`filename`))
	scanresult = append(scanresult, result.GetString(`md5`))
	scanresult = append(scanresult, result.GetString(`description`))
	return scanresult
}

// 百度webdir+ 扫描
func WebDirScan(scanfile string) []string {
	url := "https://scanner.baidu.com/enqueue"
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()
		part, err := m.CreateFormFile("archive", "scanfile.zip")
		if err != nil {
			return
		}
		file, err := os.Open(scanfile)
		if err != nil {
			return
		}
		defer file.Close()
		if _, err = io.Copy(part, file); err != nil {
			return
		}
	}()
	a, _ := http.Post(url, m.FormDataContentType(), r)
	rspBody, _ := ioutil.ReadAll(a.Body)
	result := viper.New()
	result.SetConfigType("json")
	result.ReadConfig(bytes.NewBuffer(rspBody))
	resulturl := result.GetString("url")
	scanresult := bdresult(resulturl)
	return scanresult
}

// 百度webshell结果爬取
func bdresult(url string) []string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	// 请求头设置 设置长连接
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.0.0 Safari/537.36")
	req.Header.Set("Connection", "keep-alive")
resetcheck:
	// 开始请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Http get err:", err)
	}
	//if resp.StatusCode != 200 {
	//	fmt.Println("status code:", err)
	//}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var scanresult []string

	if checkstatus := strings.Contains(string(body), "pending"); checkstatus {
		time.Sleep(3 * time.Second)
		goto resetcheck
	} else {
		fmt.Printf("百度webdir+扫描完成")
	}
	result := viper.New()
	result.SetConfigType("json")

	html := string(body)
	html = html[1 : len(html)-2]
	result.ReadConfig(bytes.NewBuffer([]byte(html)))
	if err := result.ReadConfig(bytes.NewBuffer([]byte(html))); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("找不到配置文件..")
		} else {
			fmt.Println("配置文件出错..")
		}
	}
	scanresult[1] = result.GetString(`md5`)
	//fmt.Printf(result.GetString(`data.descr`))
	//fmt.Printf(result.GetString(`md5`))
	re_json := strings.Split(html, "],")
	re_json1 := re_json[0]
	results := re_json1[9:len(re_json1)]
	// path descr
	result.ReadConfig(bytes.NewBuffer([]byte(results)))
	if err := result.ReadConfig(bytes.NewBuffer([]byte(results))); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("找不到配置文件..")
		} else {
			fmt.Println("配置文件出错..")
		}
	}
	scanresult[0] = result.GetString(`path`)
	scanresult[2] = result.GetString(`descr`)
	return scanresult
}

// Aliscan 阿里伏魔
func Aliscan() {
	// 菜狗不会js逆向啊!!!
	// 写不来 有师傅会写的可以来填坑
	// 阿里伏魔平台 https://ti.aliyun.com/#/webshell
	// 师傅们有其他的webshell在线检测网站可以多多推荐
	// 有想法一起完善这个工具的师傅联系我
	// Mail: root@tx1ee.cc
}

// Sechm 安全检测
var Sechm []string
var Secbd []string
var Secresult []string

func Seccheck(samplespath string) []string {
	if Gpushconf[5] == "true" {
		Sechm = Hmscan(samplespath)
		Secbd = WebDirScan(samplespath)
		Secresult = append(Secresult, Sechm[4])
		Secresult = append(Secresult, Secbd[3])
	} else {
		return nil
	}
	// 返回河马和百度的检查结果
	return Secresult
}
