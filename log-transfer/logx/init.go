package logx

import (
	"io"
	"log"
	"os"
)

var Log *log.Logger

func Init() {
	f, err := os.OpenFile(".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic("打开日志文件失败！" + err.Error())
	}
	Log = log.New(io.MultiWriter(f, os.Stdout), "[log-agent] ", log.LstdFlags)
}
