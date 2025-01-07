// logger/zap_lumberjack.go
package logger

import (
	"github.com/levensspel/go-gin-template/helper"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Info(msg string, function helper.FunctionCaller, data ...interface{})
	Error(msg string, function helper.FunctionCaller, data ...interface{})
	Debug(msg string, function helper.FunctionCaller, data ...interface{})
	Warn(msg string, function helper.FunctionCaller, data ...interface{})
}

type LogHandler struct {
	logger *zap.SugaredLogger
}

func NewlogHandler() *LogHandler {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    10, // Max megabytes before log is rotated
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	}

	writeSyncer := zapcore.AddSync(lumberjackLogger)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		zapcore.InfoLevel,
	)

	logger := zap.New(core, zap.AddCaller())
	return &LogHandler{
		logger: logger.Sugar(),
	}
}

func NewlogHandlerInject(i do.Injector) (LogHandler, error) {
	logger := NewlogHandler()
	return *logger, nil
}

func (l *LogHandler) Info(msg string, function helper.FunctionCaller, data ...interface{}) {
	l.logger.Infow(msg,
		"called_by", function,
		"data", data,
	)
}

func (l *LogHandler) Error(msg string, function helper.FunctionCaller, data ...interface{}) {
	l.logger.Errorw(msg,
		"called_by", function,
		"data", data,
	)
}

func (l *LogHandler) Debug(msg string, function helper.FunctionCaller, data ...interface{}) {
	l.logger.Debugw(msg,
		"called_by", function,
		"data", data,
	)
}

func (l *LogHandler) Warn(msg string, function helper.FunctionCaller, data ...interface{}) {
	l.logger.Warnw(msg,
		"called_by", function,
		"data", data,
	)
}
