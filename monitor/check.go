package monitor

import (
	"bytes"
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
			Error.Println("URL:" + url + "找不到配置文件..")
		} else {
			Error.Println("URL:" + url + "配置文件出错..")
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
		Error.Println("Http get err:", err)
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
		Info.Printf("河马webshell扫描完成" + " URL:" + url)
	}
	//fmt.Printf(string(body))
	html := strings.Replace(string(body), "\n", "", -1)
	re_json := regexp.MustCompile(`<script id="__NEXT_DATA__" type="application/json">(.*?)</script>`)
	result_json := re_json.FindString(html)
	// 获取结果JSON
	result_json = result_json[51 : len(result_json)-9]
	result := viper.New()
	// 设置配置文件类型为json 这里很重要 不是指的话会读取不到 踩坑踩坑
	result.SetConfigType("json")
	result.ReadConfig(bytes.NewBuffer([]byte(result_json)))
	if err := result.ReadConfig(bytes.NewBuffer([]byte(result_json))); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			Error.Println("URL:" + url + "找不到配置文件..")
		} else {
			Error.Println("URL:" + url + "配置文件出错..")
		}
	}
	// 返回字符数组
	var scanresult []string
	scanresult = strings.Split(result.GetString(`props.pageProps.data.summary`), "/")
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
		Error.Println("Http get err:", err)
	}
	//if resp.StatusCode != 200 {
	//	fmt.Println("status code:", err)
	//}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if checkstatus := strings.Contains(string(body), "pending"); checkstatus {
		time.Sleep(3 * time.Second)
		goto resetcheck
	} else {
		Info.Printf("百度webdir+扫描完成" + " URL:" + url)
	}
	result := viper.New()
	result.SetConfigType("json")

	html := string(body)
	html = html[1 : len(html)-2]
	result.ReadConfig(bytes.NewBuffer([]byte(html)))
	if err := result.ReadConfig(bytes.NewBuffer([]byte(html))); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			Error.Println("找不到配置文件..")
		} else {
			Error.Println("配置文件出错..")
		}
	}
	// 返回字符
	var scanresult []string
	scanresult = append(scanresult, result.GetString(`detected`))
	scanresult = append(scanresult, "0")
	scanresult = append(scanresult, result.GetString(`scanned`))
	return scanresult
}

// Sechm 安全检测
var Sechm []string
var Secbd []string
var Secresult []string

func Seccheck(samplespath string) ([]string, []string) {
	if "true" == "true" {
		Sechm = Hmscan(samplespath)
		Secbd = WebDirScan(samplespath)
	}
	// 返回河马和百度的检查结果
	return Sechm, Secbd
}
