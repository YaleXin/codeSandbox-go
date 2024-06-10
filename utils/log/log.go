package log

import (
	"bytes"
	"codeSandbox/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

type MyFormatter struct{}

// 配置日志切割
// LogFileCut 日志文件切割
func LogFileCut(fileName string) *rotatelogs.RotateLogs {
	logier, err := rotatelogs.New(
		// 切割后日志文件名称
		fileName,
		rotatelogs.WithMaxAge(30*24*time.Hour),    // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)

	if err != nil {
		panic(err)
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel:  logier,
		logrus.FatalLevel: logier,
		logrus.DebugLevel: logier,
		logrus.WarnLevel:  logier,
		logrus.ErrorLevel: logier,
		logrus.PanicLevel: logier,
	},
		// 设置分割日志样式
		&MyFormatter{})
	logrus.AddHook(lfHook)
	return logier
}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string

	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("[codeSandbox-app] [%s] [%s] [%s:%d] [msg=%s]\n",
			timestamp, entry.Level, fName, entry.Caller.Line, entry.Message)
	} else {
		newLog = fmt.Sprintf("[codeSandbox-app] [%s] [%s] [msg=%s]\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}

func ConfigLog() {
	serverConf := utils.Config.Server
	appMode := serverConf.AppMode

	os.MkdirAll("log", 0755)
	logrus.SetReportCaller(true)
	// 设置日志输出控制台样式
	logrus.SetFormatter(&MyFormatter{})
	// 按天分割
	var logFileName string
	if appMode == "dev" {
		logFileName = path.Join("log", "codeSandbox") + ".%Y%m%d_dev.log"
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logFileName = path.Join("log", "codeSandbox") + ".%Y%m%d_prod.log"
		logrus.SetLevel(logrus.InfoLevel)
	}
	// 配置日志分割
	logFileCut := LogFileCut(logFileName)
	writers := []io.Writer{
		logFileCut,
		os.Stdout}

	// 输出到控制台，方便定位到那个文件
	fileAndStdoutWriter := io.MultiWriter(writers...)
	gin.DefaultWriter = fileAndStdoutWriter

	logrus.Info("init log end")
}
