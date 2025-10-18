package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Log          *slog.Logger
	fileLogger   *lumberjack.Logger
	rotateSignal chan os.Signal
)

func InitWithFile(logPath string) error {
	if logPath == "" {
		// No file logging, keep stdout only
		fmt.Fprintf(os.Stdout, "file logging disabled\n")
		return nil
	}

	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// Setup lumberjack for rotating file logs
	fileLogger = &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}

	// Write to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, fileLogger)

	// Replace logger with file output
	Log = slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(Log)

	// Setup signal handler for manual rotation
	setupRotationSignal()

	slog.Info("file logging initialized", "path", logPath)
	return nil
}

func setupRotationSignal() {
	if fileLogger == nil {
		return
	}

	rotateSignal = make(chan os.Signal, 1)
	signal.Notify(rotateSignal, syscall.SIGHUP)

	go func() {
		for {
			<-rotateSignal
			slog.Info("manual log rotation triggered via SIGHUP, rotating now")
			if err := fileLogger.Rotate(); err != nil {
				slog.Error("failed to rotate log", "error", err)
			} else {
				slog.Info("log file rotated successfully")
			}
		}
	}()
}

// Close closes the file logger gracefully
func Close() error {
	if fileLogger != nil {
		if rotateSignal != nil {
			signal.Stop(rotateSignal)
			close(rotateSignal)
		}
		return fileLogger.Close()
	}
	return nil
}
