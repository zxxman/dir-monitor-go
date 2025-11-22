package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// LogLevel Log level type
type LogLevel int

// Log levels
const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

const (
	DefaultLogCallerDepth = 2
	DefaultLogFilePerm    = 0644
	DefaultLogDirPerm     = 0755
	LogTimeFormat         = "2006-01-02 15:04:05"
	LogBackupTimeFormat   = "20060102_150405"
)

type LoggerConfig struct {
	Level      LogLevel
	Output     io.Writer
	FilePath   string
	MaxSize    int64
	ShowCaller bool
}

// Logger Logger structure
type Logger struct {
	config LoggerConfig
	mu     sync.RWMutex
	file   *os.File
	caller int // 调用栈深度
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		err := l.file.Close()
		l.file = nil
		l.config.Output = os.Stdout
		return err
	}
	return nil
}

func NewLogger(level LogLevel, output io.Writer) *Logger {
	config := LoggerConfig{
		Level:      level,
		Output:     output,
		ShowCaller: false,
	}

	return &Logger{
		config: config,
		caller: DefaultLogCallerDepth,
	}
}

func NewLoggerWithConfig(config LoggerConfig) (*Logger, error) {
	logger := &Logger{
		config: config,
		caller: DefaultLogCallerDepth,
	}

	if config.FilePath != "" {
		err := logger.initFileLogger()
		if err != nil {
			return nil, err
		}
	}

	return logger, nil
}

func NewFileLogger(level LogLevel, filePath string, maxSize int64) (*Logger, error) {
	config := LoggerConfig{
		Level:      level,
		FilePath:   filePath,
		MaxSize:    maxSize,
		ShowCaller: false,
	}

	return NewLoggerWithConfig(config)
}

func (l *Logger) initFileLogger() error {
	dir := filepath.Dir(l.config.FilePath)
	if err := os.MkdirAll(dir, DefaultLogDirPerm); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	file, err := os.OpenFile(l.config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, DefaultLogFilePerm)
	if err != nil {
		return err
	}

	l.file = file
	l.config.Output = file
	return nil
}

func (l *Logger) SetCaller(enable bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.ShowCaller = enable
}

// Debug Log debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info Log info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn Log warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error Log error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.config.Level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	var message string
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	} else {
		message = format
	}

	l.writeTextLog(level, message)

	if l.file != nil && l.config.MaxSize > 0 {
		l.rotateLog()
	}
}

func (l *Logger) writeTextLog(level LogLevel, message string) {
	timestamp := time.Now().Format(LogTimeFormat)
	levelStr := l.levelToString(level)

	var callerInfo string
	if l.config.ShowCaller {
		_, file, line, ok := runtime.Caller(l.caller)
		if ok {
			callerInfo = fmt.Sprintf(" [%s:%d]", filepath.Base(file), line)
		}
	}

	var logEntry string
	if callerInfo != "" {
		logEntry = fmt.Sprintf("[%s] %s%s %s\n",
			timestamp, levelStr, callerInfo, message)
	} else {
		logEntry = fmt.Sprintf("[%s] %s %s\n",
			timestamp, levelStr, message)
	}

	l.config.Output.Write([]byte(logEntry))
}

// levelToString Convert log level to string
func (l *Logger) levelToString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKN"
	}
}

// rotateLog Rotate log file if needed
func (l *Logger) rotateLog() {
	if l.file == nil {
		return
	}

	info, err := l.file.Stat()
	if err != nil {
		return
	}

	if info.Size() < l.config.MaxSize {
		return
	}

	// Close current file
	l.file.Close()

	// Create backup file name
	timestamp := time.Now().Format(LogBackupTimeFormat)
	backupPath := l.config.FilePath + "." + timestamp + ".bak"

	// Rename current file to backup
	os.Rename(l.config.FilePath, backupPath)

	// Create new log file
	err = l.initFileLogger()
	if err != nil {
		// Try to reopen original file
		l.file, _ = os.OpenFile(l.config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, DefaultLogFilePerm)
		return
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.Level = level
}

func (l *Logger) GetLevel() LogLevel {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.config.Level
}

func (l *Logger) WithModule(module string) *ModuleLogger {
	return &ModuleLogger{
		logger: l,
		module: module,
	}
}

type ModuleLogger struct {
	logger *Logger
	module string
}

// Debug Log debug message with module
func (ml *ModuleLogger) Debug(format string, args ...interface{}) {
	ml.logWithModule(DEBUG, format, args...)
}

// Info Log info message with module
func (ml *ModuleLogger) Info(format string, args ...interface{}) {
	ml.logWithModule(INFO, format, args...)
}

// Warn Log warning message with module
func (ml *ModuleLogger) Warn(format string, args ...interface{}) {
	ml.logWithModule(WARN, format, args...)
}

// Error Log error message with module
func (ml *ModuleLogger) Error(format string, args ...interface{}) {
	ml.logWithModule(ERROR, format, args...)
}

func (ml *ModuleLogger) logWithModule(level LogLevel, format string, args ...interface{}) {
	if level < ml.logger.config.Level {
		return
	}

	ml.logger.mu.Lock()
	defer ml.logger.mu.Unlock()

	var message string
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	} else {
		message = format
	}

	ml.writeTextLogWithModule(level, message)

	if ml.logger.file != nil && ml.logger.config.MaxSize > 0 {
		ml.logger.rotateLog()
	}
}

func (ml *ModuleLogger) writeTextLogWithModule(level LogLevel, message string) {
	timestamp := time.Now().Format(LogTimeFormat)
	levelStr := ml.logger.levelToString(level)

	var callerInfo string
	if ml.logger.config.ShowCaller {
		_, file, line, ok := runtime.Caller(DefaultLogCallerDepth)
		if ok {
			callerInfo = fmt.Sprintf(" [%s:%d]", filepath.Base(file), line)
		}
	}

	var logEntry string
	if callerInfo != "" {
		logEntry = fmt.Sprintf("[%s] %s [%s]%s %s\n",
			timestamp, levelStr, ml.module, callerInfo, message)
	} else {
		logEntry = fmt.Sprintf("[%s] %s [%s] %s\n",
			timestamp, levelStr, ml.module, message)
	}

	ml.logger.config.Output.Write([]byte(logEntry))
}
