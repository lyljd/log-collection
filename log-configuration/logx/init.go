package logx

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var Log *log.Logger

func Init() {
	filePath := "log/.log"
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		panic("初始化日志目录失败！" + err.Error())
	}
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic("打开日志文件失败！" + err.Error())
	}
	Log = log.New(io.MultiWriter(f, os.Stdout), "[log-configuration] ", log.LstdFlags)
}
