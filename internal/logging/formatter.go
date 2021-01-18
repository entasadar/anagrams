package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

var log = logrus.New()

type AppLogger struct {
	loggerName string
}

func (l AppLogger) Error(text string) {
	log.WithField("source", l.loggerName).Error(text)
}
func (l AppLogger) Info(text string) {
	log.WithField("source", l.loggerName).Info(text)
}
func (l AppLogger) Debug(text string) {
	log.WithField("source", l.loggerName).Debug(text)
}
func (l AppLogger) Warn(text string) {
	log.WithField("source", l.loggerName).Warn(text)
}

func InitLogger(logFile string, loggingLevel logrus.Level) error {
	log.Level = loggingLevel
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	if logFile == "" {
		logFile = fmt.Sprintf("%s.log", time.Now().Format("2006-01-02 15:04:05"))
	}
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Join(workDir, "logs"), os.ModePerm); err != nil {
		return err
	}
	logPath := filepath.Join(workDir, "logs", logFile)

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		return err
	}

	return nil
}

func GetLogger(name string) AppLogger {
	return AppLogger{loggerName: name}
}
