package config

import (
	logger "github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
	"io"
	"os"
)

func SetupLogger() *os.File {
	logger.SetFormatter(&ecslogrus.Formatter{})
	logger.SetLevel(logger.TraceLevel)
	logFilePath := "./logs/out.log"
	if err := os.MkdirAll("./logs", os.ModePerm); err != nil {
		logger.Fatal(err, "Failed to create log directory")
	}
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatal(err, "Failed to open the log file")
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	logger.SetOutput(multiWriter)
	return file
}
