package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// getCallerInfo returns full/relative path, line, and function name
func getCallerInfo(skip int) (string, int, string) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "???", 0, "unknown"
	}

	fn := runtime.FuncForPC(pc)
	funcName := "unknown"
	if fn != nil {
		funcName = fn.Name()
		if idx := strings.LastIndex(funcName, "/"); idx != -1 {
			funcName = funcName[idx+1:]
		}
	}

	wd, err := os.Getwd()
	if err == nil {
		if rel, err := filepath.Rel(wd, file); err == nil {
			file = rel
		}
	}

	return file, line, funcName
}

// Debug prints debug logs with caller file, line, and function name
func Debug(format string, a ...interface{}) {
	file, line, funcName := getCallerInfo(2)
	timestamp := time.Now().Format("15:04:05")
	msg := fmt.Sprintf(format, a...)

	fmt.Printf("\x1b[36m[%s] [DEBUG]\x1b[0m \x1b[90m%s:%d (%s)\x1b[0m | %s\n",
		timestamp,
		file,
		line,
		funcName,
		msg,
	)
}

// LogErr prints error logs with caller file, line, and function name
func LogErr(format string, a ...interface{}) {
	file, line, funcName := getCallerInfo(2)
	timestamp := time.Now().Format("15:04:05")
	msg := fmt.Sprintf(format, a...)

	fmt.Printf("\x1b[31m[%s] [ERROR]\x1b[0m \x1b[90m%s:%d (%s)\x1b[0m | %s\n",
		timestamp,
		file,
		line,
		funcName,
		msg,
	)
}

// LogInfo prints general information logs
func LogInfo(format string, a ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("\x1b[32m[%s] [INFO]\x1b[0m | %s\n", timestamp, msg)
}

// LogWarn prints warning logs
func LogWarn(format string, a ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("\x1b[38;5;214m[%s] [Warning]\x1b[0m | %s\n", timestamp, msg)
}
