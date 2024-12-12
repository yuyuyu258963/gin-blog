package file

import (
	"io"
	"mime/multipart"
	"os"
	"path"
)

// 获取文件大小
func GetSize(f multipart.File) (int, error) {
	content, err := io.ReadAll(f)
	if err != nil {
		return 0, err
	}
	return len(content), err
}

// 获取文件格式
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

// 检查文件是否存在
func CheckExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

// 检查是否有权限问题
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

// 如果不存在则新建文件夹
func IsNotExistMkDir(src string) error {
	if notExist := CheckExist(src); notExist {
		if err := MkDir(src); err != nil {
			return err
		}
	}
	return nil
}

// 新建文件夹，如果有父路径不存在的话则循环创建
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}
	return err
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}
