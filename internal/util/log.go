package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"io"
	"os"
)

var Logger *zap.SugaredLogger

// InitLog 初始化日志
func InitLog() {
	writeSyncer := GetLogWriter()
	encoder := GetEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	Logger = zap.New(core).Sugar()
}
func GetEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
func GetLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create("./k8s_install.logger")
	//日志同时输出到文件和终端
	w := io.MultiWriter(os.Stdout, file)
	return zapcore.AddSync(w)
}
