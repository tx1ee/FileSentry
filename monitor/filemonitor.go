package monitor

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gookit/color"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// FileExt 文件后缀获取
/*
调用测试代码
FileExts := []string{".php", ".aspx"}
a1 := FileExt("aaaa.php")
if Strinarr(a1[1], FileExts) {
	fmt.Println("ok")
}
*/

// 事件ID

func FileExt(filename string) []string {
	//获取文件名称带后缀
	fileNameWithSuffix := path.Base(filename)
	//获取文件的后缀(文件类型)
	fileType := path.Ext(fileNameWithSuffix)
	//获取文件名称(不带后缀)
	fileNameOnly := strings.TrimSuffix(fileNameWithSuffix, fileType)
	return []string{fileNameOnly, fileType}
}

// Strinarr 判断文件后缀是否在列表
func Strinarr(target string, str_array []string) bool {
	for _, element := range str_array {
		if target == element {
			return true
		}
	}
	return false
}

// NotifyFile 参考csdn
// https://blog.csdn.net/finghting321/article/details/102852746?spm=1001.2101.3001.6650.2&utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7Edefault-2-102852746-blog-121683205.pc_relevant_default&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7Edefault-2-102852746-blog-121683205.pc_relevant_default&utm_relevant_index=4
type NotifyFile struct {
	watch *fsnotify.Watcher
}

func NewNotifyFile() *NotifyFile {
	w := new(NotifyFile)
	w.watch, _ = fsnotify.NewWatcher()
	return w
}

func (this *NotifyFile) WatchDir(dir string, Edirs []string, bfilesuffix []string) {
	//func (this *NotifyFile) WatchDir(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//判断是否为目录，监控目录,目录下文件也在监控范围内，不需要加
		if info.IsDir() {
			path, err := filepath.Abs(path)
			// 判断文件夹是否在排除
			if Strinarr(path, Edirs) {
				this.watch.Remove(path)
				color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), "  删除监控:", path, "\n")
				Info.Println("删除监控:", path)
			} else {
				if err != nil {
					return err
				}
				err = this.watch.Add(path)
				if err != nil {
					return err
				}
				color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), "  添加监控:", path, "\n")
				Info.Println("添加监控:", path)
			}

		}
		return nil
	})
	go this.WatchEvent(bfilesuffix) //协程
}

// WatchEvent 文件监控事件监听
func (this *NotifyFile) WatchEvent(bfilesuffix []string) {
	//var EventID string = RandString(8)
	//// 样本压缩文件名字
	//var zipname string = ".\\samples\\" + time.Now().Format("2006-01-02-(") + EventID + ").zip"
	for {
		select {
		case ev := <-this.watch.Events:
			{
				// 文件\文件创建 新增文件夹添加监控
				if ev.Op&fsnotify.Create == fsnotify.Create {
					//获取新创建文件的信息，如果是目录，则加入监控中
					file, err := os.Stat(ev.Name)
					// 判断文件是否为空
					if err == nil && file.IsDir() {
						this.watch.Add(ev.Name)
						color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), "  添加监控:", ev.Name, "\n")
						color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), "  创建文件夹:", ev.Name, "\n")
						Info.Println("添加监控: ", ev.Name)
						Info.Println("创建文件夹: ", ev.Name)
					} else {
						// 判断文件是否为空
						filenotnull, _ := ioutil.ReadFile(ev.Name)
						//fmt.Println(filenotnull)
						// 事件ID
						EventID := RandString(8)
						// 样本压缩文件名字
						zipname := ".\\samples\\" + time.Now().Format("2006-01-02-(") + EventID + ").zip"
						// 判断创建文件后缀是否在列表内
						fileType := FileExt(ev.Name)
						// 判断文件后缀是否在列表
						if Strinarr(fileType[1], bfilesuffix) && len(filenotnull) > 0 {
							color.Redp("[WARNING] ", time.Now().Format("2006-01-02 15:04:05"), "  敏感文件后缀创建:", ev.Name, "\n")
							Warning.Println("敏感文件后缀创建:", ev.Name)
							// 消息模板
							Dingpushinfo[0] = EventID
							Dingpushinfo[1] = "敏感文件后缀创建"
							Dingpushinfo[2] = ev.Name
							Dingpush(Gpushconf[0], Gpushconf[1], Dingpushinfo)
							// 样本文件压缩
							Enzip(ev.Name, zipname, fileType[1], false)
							fmt.Println(zipname)
							// webshell安全检测
							//secresult := Seccheck(zipname)
							//fmt.Println(secresult)
							//fmt.Println(EventID)
							//if secresult == nil {
							//	Dingpushinfo[3] = secresult[0]
							//	Dingpushinfo[4] = secresult[1]
							//}
						} else {
							color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), "  创建文件:", ev.Name, "\n")
							Info.Println("创建文件:", ev.Name)
						}
					}
				}
			}
			// 文件\文件夹修改
			if ev.Op&fsnotify.Write == fsnotify.Write {
				fi, err := os.Stat(ev.Name)
				if err == nil && fi.IsDir() {
					// 文件夹更新看着很乱
					//fmt.Println("文件夹更新 : ", ev.Name)
				} else {
					// 文件是否为空
					filenotnull, _ := ioutil.ReadFile(ev.Name)
					// 事件ID
					//EventID := RandString(8)
					// 样本压缩文件名字
					//zipname := ".\\samples\\" + time.Now().Format("2006-01-02-(") + EventID + ").zip"
					fileType := FileExt(ev.Name)
					if Strinarr(fileType[1], bfilesuffix) && len(filenotnull) > 0 {
						// 判断webshll检测是否开启 开启则上传到河马、webdir+检测是否为木马
						color.Redp("[WARNING] ", time.Now().Format("2006-01-02 15:04:05"), "  敏感文件写入:", ev.Name, "\n")
						Warning.Println("敏感文件写入:", ev.Name)
						//time.Sleep(time.Second * 10)
						Dingpushinfo[4] = "敏感文件写入"
						Dingpushinfo[5] = ev.Name
						Dingpush(Gpushconf[0], Gpushconf[1], Dingpushinfo)
						// 样本文件压缩
						//Enzip(ev.Name, zipname, fileType[1], false)
						// webshell安全检测
						//secresult := Seccheck(zipname)
						//fmt.Println(EventID)
						//if Gpushconf[5] != "true" {
						//	fmt.Println("未开启安全检测功能")
						//}

					} else {
						color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), "  写入文件:", ev.Name, "\n")
						Info.Println("写入文件:", ev.Name)
					}
				}
			}
			// 文件\文件夹删除
			if ev.Op&fsnotify.Remove == fsnotify.Remove {
				//如果删除文件是目录，则移除监控
				fi, err := os.Stat(ev.Name)
				if err == nil && fi.IsDir() {
					this.watch.Remove(ev.Name)
					color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), "  删除监控:", ev.Name, "\n")
					color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), "  删除文件夹:", ev.Name, "\n")
					Info.Println("删除监控:", ev.Name)
					Info.Println("删除文件夹:", ev.Name)
				} else {
					color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), "  删除文件:", ev.Name, "\n")
					Info.Println("删除文件:", ev.Name)
				}
			}
			// 文件\文件夹更新
			if ev.Op&fsnotify.Rename == fsnotify.Rename {
				color.Greenp("[INFO] ", time.Now().Format("2006-01-02 15:04:05"), "  重命名文件:", ev.Name, "\n")
				Info.Println("重命名文件:", ev.Name)
				this.watch.Remove(ev.Name)
			}

		case err := <-this.watch.Errors:
			{
				Info.Println("error : ", err)
				return
			}
		}
	}
}
