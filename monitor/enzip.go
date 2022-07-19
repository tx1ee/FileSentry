package monitor

import (
	"bytes"
	"github.com/alexmullins/zip"
	"io"
	"os"
)

// Enzip 压缩包文件
func Enzip(webshellpath string, outputpzip string, Filetype string, ispwden bool) {
	// 木马文件读取
	var w io.Writer
	contents, err := os.ReadFile(webshellpath)
	if err != nil {
		panic(err)
	}
	fzip, err := os.Create(outputpzip)
	if err != nil {
		//Warning.Println(err)
	}
	zipw := zip.NewWriter(fzip)
	defer zipw.Close()
	// 文件压缩判断是否加密
	if ispwden {
		w, err = zipw.Encrypt(`shellfile`+Filetype, `Sentry@123@!#`)
		if err != nil {
			Warning.Println(err)
		}
	} else {
		w, err = zipw.Create(`shellfile` + Filetype)
		if err != nil {
			Warning.Println(err)
		}
	}
	_, err = io.Copy(w, bytes.NewReader(contents))
	if err != nil {
		Warning.Println(err)
	} else {
		Info.Println(outputpzip, "文件压缩完完成")
	}
	zipw.Flush()
}
