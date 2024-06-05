package util

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"os"
	"server.com/daren/config"
	"time"
)

func InitLogger() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	if config.IsLocalDev() {
		logrus.SetOutput(os.Stdout)
	} else {
		logrus.SetOutput(ioutil.Discard)
	}

	// Only logrus the warning severity or above.
	if config.IsDev() {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logPath := config.GetLogDir()
	logMain, err := rotatelogs.New(
		logPath+"main.log.%Y%m%d",
		rotatelogs.WithMaxAge(30*24*time.Hour),
	)
	if err != nil {
		log.Fatal("初始化日志失败")
	}

	logError, err := rotatelogs.New(
		logPath+"error.log.%Y%m%d",
		rotatelogs.WithMaxAge(30*24*time.Hour),
	)
	if err != nil {
		log.Fatal("初始化日志失败")
	}

	var formatter logrus.Formatter
	formatter = new(logrus.JSONFormatter)
	//为不同级别设置不同的输出目的
	lfHook := NewHook(
		WriterMap{
			logrus.DebugLevel: logMain,
			logrus.InfoLevel:  logMain,
			logrus.WarnLevel:  logMain,
			logrus.ErrorLevel: logError,
			logrus.FatalLevel: logError,
			logrus.PanicLevel: logError,
		},
		formatter,
	)
	logrus.AddHook(lfHook)
}
