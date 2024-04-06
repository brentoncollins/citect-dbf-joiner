package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"sync"
)

var (
	loggerInstance *logrus.Logger
	once           sync.Once
)

func Logger(logFileLocation string) *logrus.Logger {
	once.Do(func() {
		loggerInstance = logrus.New()
		loggerInstance.SetLevel(logrus.InfoLevel)

		logFile, err := os.OpenFile(logFileLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Printf("Failed to open log file: %v", err)
		} else {
			loggerInstance.SetOutput(io.MultiWriter(os.Stdout, logFile))
		}
	})

	return loggerInstance
}
