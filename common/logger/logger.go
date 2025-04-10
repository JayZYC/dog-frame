package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path"
	"runtime"
	"sync"
)

var (
	f    *facade
	once sync.Once
)

type facade struct {
	_logger *zap.Logger
}

func Debug(ctx context.Context, msg string, kv ...interface{}) {
	logFacade().log(ctx, zapcore.DebugLevel, msg, kv...)
}

func Info(ctx context.Context, msg string, kv ...interface{}) {
	logFacade().log(ctx, zapcore.InfoLevel, msg, kv...)
}

func Warn(ctx context.Context, msg string, kv ...interface{}) {
	logFacade().log(ctx, zapcore.WarnLevel, msg, kv...)
}

func Error(ctx context.Context, msg string, kv ...interface{}) {
	logFacade().log(ctx, zapcore.ErrorLevel, msg, kv...)
}

func (f *facade) log(ctx context.Context, lvl zapcore.Level, msg string, kv ...interface{}) {
	fields := makeLogFields(ctx, kv...)
	ce := f._logger.Check(lvl, msg)
	ce.Write(fields...)
}

func logFacade() *facade {
	// 使用sync.Once确保初始化代码只执行一次
	once.Do(func() {
		f = &facade{
			_logger: _logger,
		}
	})
	return f
}

// 组装日志行中的字段
func makeLogFields(ctx context.Context, kv ...interface{}) []zap.Field {
	// 保证日志信息以键值对的形式成对出现
	if len(kv)%2 != 0 {
		kv = append(kv, "unknown")
	}
	// 设置日志行信息中的追踪字段
	kv = append(kv, "traceid", ctx.Value("traceid"), "spanid", ctx.Value("spanid"), "pspanid", ctx.Value("pspanid"))
	// 增加日志调用者信息
	funcName, file, line := getLoggerCallerInfo()
	kv = append(kv, "func", funcName, "file", file, "line", line)
	fields := make([]zap.Field, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		k := fmt.Sprintf("%v", kv[i])
		fields = append(fields, zap.Any(k, kv[i+1]))
	}

	return fields
}

// getLoggerCallerInfo 日志调用者信息 -- 方法名, 文件名, 行号
func getLoggerCallerInfo() (funcName, file string, line int) {
	pc, file, line, ok := runtime.Caller(4) // 回溯拿调用日志方法的业务函数的信息
	if !ok {
		return
	}
	file = path.Base(file)
	funcName = runtime.FuncForPC(pc).Name()
	return
}
