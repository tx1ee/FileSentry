package monitor

import (
	"io"
	"log"
	"os"
	"time"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	dir, _ := os.Getwd()
	nowtime := time.Now().Format("2006-01-02")
	f, err := os.OpenFile(dir+"\\logs\\"+nowtime+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	Trace = log.New(io.MultiWriter(f), "[TRACE] ", log.Ldate|log.Ltime)
	Info = log.New(io.MultiWriter(f), "[INFO] ", log.Ldate|log.Ltime)
	Warning = log.New(io.MultiWriter(f), "[WARNING] ", log.Ldate|log.Ltime)
	Error = log.New(io.MultiWriter(f), "[ERROR] ", log.Ldate|log.Ltime)
}
