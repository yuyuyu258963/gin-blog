package file

import (
	"fmt"
	"io"
	"log"
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
func CheckNotExist(src string) bool {
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
	if notExist := CheckNotExist(src); notExist {
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

// MustOpen: maximize trying to open file
func MustOpen(fileName string, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err:%v", err)
	}

	src := path.Join(dir, filePath)
	perm := CheckPermission(src)
	if perm {
		return nil, fmt.Errorf("CheckPermission Permission denied")
	}

	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("IsNotExistMkDir sec:%s, err:%v", src, err)
	}

	log.Println("trying to open... ", path.Join(src, fileName))
	f, err := Open(path.Join(src, fileName), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Open err src:%v, fileName:%v, err:%v", src, fileName, err)
	}
	return f, nil
}
