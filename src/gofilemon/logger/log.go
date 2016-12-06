package logger

import (
	"fmt"
	"os"
	"path/filepath"
)

var logFileFd *os.File

func GetLogFD() (*os.File, error) {
	if logFileFd != nil {
		return logFileFd, nil
	}
	fileName := filepath.Join("/tmp", fmt.Sprintf("%s.log", filepath.Base(os.Args[0])))
	fd, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	logFileFd = fd
	return logFileFd, nil
}
