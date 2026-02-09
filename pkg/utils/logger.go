package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// LogLevel 日志级别
type LogLevel int

const (
	// DebugLevel 调试级别
	DebugLevel LogLevel = iota
	// InfoLevel 信息级别
	InfoLevel
	// WarnLevel 警告级别
	WarnLevel
	// ErrorLevel 错误级别
	ErrorLevel
	// FatalLevel 致命级别
	FatalLevel
)

// Logger 日志记录器
type Logger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	level       LogLevel
}

// NewLogger 创建新的日志记录器
func NewLogger(level LogLevel) *Logger {
	return &Logger{
		debugLogger: log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		fatalLogger: log.New(os.Stderr, "[FATAL] ", log.Ldate|log.Ltime|log.Lshortfile),
		level:       level,
	}
}

// GetLogLevelFromString 从字符串获取日志级别
func GetLogLevelFromString(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// Debug 记录调试级别日志
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= DebugLevel {
		l.debugLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Info 记录信息级别日志
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= InfoLevel {
		l.infoLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Warn 记录警告级别日志
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= WarnLevel {
		l.warnLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Error 记录错误级别日志
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= ErrorLevel {
		l.errorLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Fatal 记录致命级别日志并退出程序
func (l *Logger) Fatal(format string, v ...interface{}) {
	if l.level <= FatalLevel {
		l.fatalLogger.Output(2, fmt.Sprintf(format, v...))
		os.Exit(1)
	}
}

// MaskAPIKey 对API Key进行脱敏处理
func MaskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "***"
	}
	return apiKey[:4] + "***" + apiKey[len(apiKey)-4:]
}
