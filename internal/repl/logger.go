package repl

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	file "github.com/YouWantToPinch/pincher-cli/internal/filemgr"
)

// Logger handles all slog output.
// Letting this type have the final say allows for ...what? Why did I do all this again?
type Logger struct {
	logFile *os.File
	slogger *slog.Logger
}

func (l *Logger) getCurrentLog() string {
	return "log_" + time.Now().Format("2006-01-02")
}

func (l *Logger) New(level slog.Level) error {
	path, _ := file.GetLogFilepath(l.getCurrentLog())

	err := os.MkdirAll(filepath.Dir(path), 0o755)
	if err != nil {
		return fmt.Errorf("could not open or create directory: %s", err.Error())
	}

	l.logFile, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("could not open or create log file: %s", err.Error())
	}

	l.slogger = slog.New(slog.NewJSONHandler(l.logFile,
		&slog.HandlerOptions{Level: level}))
	slog.SetDefault(l.slogger)

	return nil
}

// Close ensures the removal of any empty log file created on the day that
// the session ended
func (l *Logger) Close() error {
	if l.logFile != nil {
		fileInfo, err := l.logFile.Stat()
		if err != nil {
			return fmt.Errorf("could not get file info for log: %s", err.Error())
		}
		if fileInfo.Size() == 0 {
			path, _ := file.GetLogFilepath(l.getCurrentLog())

			if err := os.Remove(path); err != nil {
				return fmt.Errorf("could not remove log file: %s", err.Error())
			}
		} else if err := l.logFile.Close(); err != nil {
			panic(err)
		}
	}

	return nil
}
