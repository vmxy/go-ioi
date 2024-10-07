package util

import (
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	// 设置日志格式
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// 日志分级
	log.Info("这是一个信息日志")
	log.Warn("这是一个警告日志")
	log.Error("这是一个错误日志")
	log.Fatal("这是一个致命错误日志") // 会调用 os.Exit(1)
	log.Panic("这是一个恐慌日志")   // 会触发 panic
}

var Log = log
