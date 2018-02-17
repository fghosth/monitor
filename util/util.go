package util

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Log   *logrus.Logger
	Viper *viper.Viper
)

func init() {
	Viper = viper.New()

	Log = logrus.New()
	// 以json格式显示.
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	Log.Formatter = &logrus.JSONFormatter{}
	viper.SetDefault("ContentDir", "content")

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	// logrus.SetOutput(os.Stdout)
	Log.Out = os.Stdout

	// file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	// if err == nil {
	// 	Log.Out = file
	// } else {
	// 	Log.Info("日志文件打开失败，使用默认输出")
	// }

	// Only log the warning severity or above.
	// logrus.SetLevel(logrus.WarnLevel)
	Log.Level = logrus.InfoLevel
}
