package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path"
	"runtime"
)

var (
	log *logger
)

type logger struct {
	_logger *zap.Logger
}

func (l *logger) log(ctx context.Context, lvl zapcore.Level, msg string, kv ...interface{}) {
	// 如果键值对漏了，那么补充一个unknown，确保日志能顺利打印
	if len(kv)%2 != 0 {
		kv = append(kv, "unknown")
	}
	// 日志信息增加追踪参数
	kv = append(kv, "traceid", ctx.Value("traceid"), "spanid", ctx.Value("spanid"), "pspanid", ctx.Value("pspanid"))
	// 增加日志调用者信息
	funcName, file, line := l.getLoggerCallerInfo()
	kv = append(kv, "func", funcName, "caller", file, "line", line)
	fields := make([]zap.Field, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		k := fmt.Sprintf("%v", kv[i])
		fields = append(fields, zap.Any(k, kv[i+1]))
	}
	ce := l._logger.Check(lvl, msg)
	ce.Write(fields...)
}

// getLoggerCallerInfo 日志调用者信息 -- 方法名, 文件名, 行号
func (l *logger) getLoggerCallerInfo() (funcName string, file string, line int) {
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return
	}
	file = path.Base(file)
	funcName = runtime.FuncForPC(pc).Name()
	return
}

func (l *logger) Debug(ctx context.Context, msg string, kv ...interface{}) {
	l.log(ctx, zapcore.DebugLevel, msg, kv...)
}

func (l *logger) Info(ctx context.Context, msg string, kv ...interface{}) {
	l.log(ctx, zapcore.InfoLevel, msg, kv...)
}

func (l *logger) Warn(ctx context.Context, msg string, kv ...interface{}) {
	l.log(ctx, zapcore.WarnLevel, msg, kv...)
}

func (l *logger) Error(ctx context.Context, msg string, kv ...interface{}) {
	l.log(ctx, zapcore.ErrorLevel, msg, kv...)
}

func Debug(ctx context.Context, msg string, kv ...interface{}) {
	log.Debug(ctx, msg, kv...)
}

func Info(ctx context.Context, msg string, kv ...interface{}) {
	log.Info(ctx, msg, kv...)
}

func Warn(ctx context.Context, msg string, kv ...interface{}) {
	log.Warn(ctx, msg, kv...)
}

func Error(ctx context.Context, msg string, kv ...interface{}) {
	log.Error(ctx, msg, kv...)
}
