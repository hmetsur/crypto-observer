// pkg/logger/logger.go
package logger

import (
	"os"
	"strings"

	"crypto-observer/pkg/config"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// parseLevel — отдельно, чтобы легко покрыть тестами.
func parseLevel(s string) logrus.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return logrus.DebugLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

func Init() {
	if log != nil {
		return
	}
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetFormatter(&logrus.JSONFormatter{})

	// 1) приоритет окружения (LOG_LEVEL), 2) конфиг, 3) INFO по умолчанию
	levelStr := strings.TrimSpace(os.Getenv("LOG_LEVEL"))
	if levelStr == "" {
		levelStr = config.C().Log.Level
	}
	l.SetLevel(parseLevel(levelStr))

	log = l
}

// L — безопасный геттер: если что-то вызвало L() до Init(), мы не упадём.
func L() *logrus.Logger {
	if log == nil {
		Init()
	}
	return log
}

type Fields = logrus.Fields
