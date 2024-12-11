package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	LogSavePath = "./runtime/logs/"
	LogSaveName = "log"
	LogFileExt  = "log"
	TimeFormat  = "20060102"
)

// 获取log文件保存的位置
func getLogFilePath() string {
	return LogSavePath
}

// 获取日志文件完整路径
func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt)

	return fmt.Sprintf("%s/%s", prefixPath, suffixPath)
}

// 常打开日志文件
func openLogFile(filePath string) *os.File {
	// 返回文件信息结构描述文件
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		mkDir()
	case os.IsPermission(err):
		log.Fatalf("Permission: %v", err)
	}

	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile : %v", err)
	}
	return handle
}

// 创建日志文件的文件夹
func mkDir() {
	dir, _ := os.Getwd()
	// os.ModePerm = 0777 即表示拥有所有权限
	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
