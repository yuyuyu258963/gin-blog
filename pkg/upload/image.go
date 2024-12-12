package upload

import (
	"fmt"
	"gin_example/pkg/file"
	"gin_example/pkg/logging"
	"gin_example/pkg/setting"
	"gin_example/pkg/util"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

// 获取图片的原始url
func GetImageFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetImagePath() + name
}

// 获取文件名，其中name会通过MD5编码
func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)

	return fileName + ext
}

// 获取文件保存的目录
func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

// 获取文件保存在本地的完整地址
func GetImageFullPath() string {
	return path.Join(setting.AppSetting.RuntimeRootPath, GetImagePath())
}

// 检查文件的后缀是否被允许
func CheckImageExt(filename string) bool {
	ext := file.GetExt(filename)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		// 大小写映射后变化，如 .jpg = .JPG
		if strings.EqualFold(allowExt, ext) {
			return true
		}
	}
	return false
}

// 检查文件的大小是否合法
func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}

	return size < setting.AppSetting.ImageMaxSize
}

// 检查图片
func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os Getwd err: %v", err)
	}

	err = file.IsNotExistMkDir(path.Join(dir, src))
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err:%v", err)
	}

	perm := file.CheckPermission(src)
	if perm {
		return fmt.Errorf("file.CheckPermission denied src:%s", src)
	}

	return nil
}
