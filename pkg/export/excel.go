package export

import (
	"gin_example/pkg/setting"
	"path"
)

// 获取文件所在的位置
func GetExcelFullUrl(name string) string {
	return path.Join(setting.AppSetting.PrefixUrl, GetExcelPath(), name)
}

// 获取文件所在的路径
func GetExcelPath() string {
	return setting.AppSetting.ExportSavePath
}

// 获取文件实际存储的相对路径
func GetExcelFullPath() string {
	return path.Join(setting.AppSetting.RuntimeRootPath, GetExcelPath())
}
