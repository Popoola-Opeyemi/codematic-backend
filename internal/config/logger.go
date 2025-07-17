package config

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerManager handles the logger instance and its cleanup
type LoggerManager struct {
	Logger     *zap.Logger
	fileWriter *os.File
	mu         sync.Mutex
	closed     bool
}

var (
	loggerManager *LoggerManager
	once          sync.Once
)

// InitLogger initializes a new logger with file and console output
func InitLogger() *LoggerManager {
	once.Do(func() {
		loggerManager = initLoggerInternal()
	})
	return loggerManager
}

// GetLogger returns the current logger instance
func GetLogger() *zap.Logger {
	if loggerManager == nil {
		return zap.NewNop()
	}
	return loggerManager.Logger
}

// initLoggerInternal creates and configures the logger
func initLoggerInternal() *LoggerManager {
	logPath := "logs/app.log"
	logDir := filepath.Dir(logPath)

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic("failed to create log directory: " + err.Error())
	}

	// Open log file
	fileWriter, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}
	fileSyncer := zapcore.AddSync(fileWriter)

	// JSON encoder for file
	jsonEncoder := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
		LineEnding:    zapcore.DefaultLineEnding,
	}
	fileCore := zapcore.NewCore(zapcore.NewJSONEncoder(jsonEncoder), fileSyncer, zap.DebugLevel)

	// Console encoder with color and custom format
	consoleEncoder := zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		CallerKey:  "caller",
		MessageKey: "msg",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("[" + t.Format(time.RFC3339) + "]")
		},
		EncodeCaller: func(c zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("[" + c.TrimmedPath() + "]")
		},
		EncodeLevel: zapcore.CapitalColorLevelEncoder,
		LineEnding:  zapcore.DefaultLineEnding,
	}
	consoleCore := zapcore.NewCore(zapcore.NewConsoleEncoder(consoleEncoder), zapcore.AddSync(os.Stdout), zap.DebugLevel)

	// Combine both cores
	core := zapcore.NewTee(fileCore, consoleCore)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &LoggerManager{
		Logger:     logger,
		fileWriter: fileWriter,
		closed:     false,
	}
}

func (lm *LoggerManager) Close() error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if lm.closed {
		return nil
	}

	// Attempt to sync only fileWriter, skip os.Stdout
	if lm.fileWriter != nil {
		if err := lm.fileWriter.Sync(); err != nil {
			os.Stderr.WriteString("failed to sync file writer: " + err.Error() + "\n")
		}
	}

	if err := lm.Logger.Sync(); err != nil {
		// Ignore known stdout/stderr sync errors
		if !isIgnorableSyncError(err) {
			os.Stderr.WriteString("failed to sync logger: " + err.Error() + "\n")
		}
	}

	if lm.fileWriter != nil {
		if err := lm.fileWriter.Close(); err != nil {
			return err
		}
		lm.fileWriter = nil
	}

	lm.closed = true
	return nil
}

func isIgnorableSyncError(err error) bool {
	// Match the common error text
	return err.Error() == "sync /dev/stdout: invalid argument" || err.Error() == "sync /dev/stderr: invalid argument"
}
