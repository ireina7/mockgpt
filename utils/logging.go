package utils

import (
	"log"
	"os"
)

func InitLogger() error {
	loggingFile := os.Getenv("log_file")
	file, err := os.OpenFile(loggingFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(file)
	return nil
}
