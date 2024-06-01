package logs

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

type LogID string

const (
	red    = "\033[31m"
	yellow = "\033[33m"
	reset  = "\033[0m"
)

func Infof(format string, v ...any) {
	file, line := caller()
	log.Printf(fmt.Sprintf("%s:%d [INFO] %s\n", file, line, format), v...)
}

func Warnf(format string, v ...any) {
	file, line := caller()
	log.Printf(fmt.Sprintf("%s%s:%d [WARN] %s%s\n", yellow, file, line, format, reset), v...)
}

func Errorf(format string, v ...any) {
	file, line := caller()
	log.Printf(fmt.Sprintf("%s%s:%d [ERROR] %s%s\n", red, file, line, format, reset), v...)
}

func Fatalf(format string, v ...any) {
	file, line := caller()
	log.Fatalf(fmt.Sprintf("%s%s:%d [FATAL] %s%s\n", red, file, line, format, reset), v...)
}

func CtxInfo(ctx context.Context, format string, v ...any) {
	file, line := caller()
	log.Printf(fmt.Sprintf("%s:%d [%v] [INFO] %s\n", file, line, ctx.Value(LogID("logID")), format), v...)
}

func CtxWarn(ctx context.Context, format string, v ...any) {
	file, line := caller()
	log.Printf(fmt.Sprintf("%s%s:%d [%v] [WARN] %s%s\n", yellow, file, line, ctx.Value(LogID("logID")), format, reset), v...)
}

func CtxError(ctx context.Context, format string, v ...any) {
	file, line := caller()
	log.Printf(fmt.Sprintf("%s%s:%d [%v] [ERROR] %s%s\n", red, file, line, ctx.Value(LogID("logID")), format, reset), v...)
}

func CtxFatal(ctx context.Context, format string, v ...any) {
	file, line := caller()
	log.Fatalf(fmt.Sprintf("%s%s:%d [%v] [FATAL] %s%s\n", red, file, line, ctx.Value(LogID("logID")), format, reset), v...)
}

func caller() (string, int) {
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Base(file)
	return file, line
}
