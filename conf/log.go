package conf

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var Log *zap.SugaredLogger

// 初始化日志库;
// dependsOn: 配置文件加载
func initRotateLog() {
	logConfig := App.Log

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = timeEncodeFunc
	encoderConfig.TimeKey = "time"
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	var writer zapcore.WriteSyncer
	var logLevel zapcore.Level
	if !logConfig.ConsoleOnly {
		// 获取当前工作目录
		pwd, err := os.Getwd()
		fmt.Println(pwd)
		if nil != err {
			panic(err)
		}

		// 创建日志目录
		os.Mkdir(logConfig.Directory, os.ModePerm)

		routateWriter, err := rotatelogs.New(
			pwd + "/" + logConfig.FilePattern,
			rotatelogs.WithLinkName(logConfig.FileLink),
			rotatelogs.WithMaxAge(24 * time.Hour * 30),
			rotatelogs.WithRotationTime(24 * time.Hour),
		)

		if nil != err {
			panic(err)
		}

		logLevel = zapcore.InfoLevel
		writer = zapcore.AddSync(routateWriter)

	} else {
		logLevel = zapcore.DebugLevel
		writer = zapcore.AddSync(os.Stdout)
	}

	logCore := zapcore.NewCore(encoder, writer, logLevel)
	// logger := zap.New(logCore, zap.AddCaller())
	logger := zap.New(logCore)
	Log = logger.Sugar()

	Log.Info("log initialized")
}

func timeEncodeFunc(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
