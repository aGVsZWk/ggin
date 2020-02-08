package logging

import (
	"time"
	"os"
	"fmt"
	"log"
)

var (
	LogSavePath = "runtime/logs/"
	LogSaveName = "log"
	LogFileExt  = "log"
	TimeFormat  = "200601012"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s", LogSavePath)
}

func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt)

	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

func mkDir(filePath string) {
	dir, _ := os.Getwd()
	//fmt.Println(filePath)		// runtime/logs/
	//fmt.Println(getLogFilePath())		// runtime/logs/
	//err := os.MkdirAll(dir + "/" + getLogFilePath(), os.ModePerm)
	err := os.MkdirAll(dir + "/" + filePath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		//fmt.Println(filePath)	// runtime/logs/log202002028.log
		mkDir(getLogFilePath())
	case os.IsPermission(err):
		log.Fatalf("Permission: %v", err)
	}
	handle, err := os.OpenFile(filePath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile : %v", err)
	}
	return handle
}
