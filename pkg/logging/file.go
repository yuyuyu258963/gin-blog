package logging

import (
	"fmt"
	"gin_example/pkg/file"
	"gin_example/pkg/setting"
	"os"
	"path"
	"time"
)

// 获取log文件保存的位置
func getLogFilePath() string {
	return setting.AppSetting.LogSavePath
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s", setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat), setting.AppSetting.LogFileExt)
}

// 获取日志文件完整路径
func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := getLogFileName()

	return fmt.Sprintf("%s/%s", prefixPath, suffixPath)
}

// 常打开日志文件
func openLogFile(filename, filePath string) (*os.File, error) {
	// 返回文件信息结构描述文件

	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err : %v", err)
	}

	src := path.Join(dir, filePath)
	perm := file.CheckPermission(src)
	if perm {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}
	err = file.IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("failed file.IsNotExistMkDir src: %s err %v", src, err)
	}

	fmt.Println("log file full path", path.Join(src, filename))

	handle, err := file.Open(path.Join(src, filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to OpenFile Err:%w", err)
	}
	return handle, nil
}
