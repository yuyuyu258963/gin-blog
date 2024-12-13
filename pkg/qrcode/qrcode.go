package qrcode

import (
	"gin_example/pkg/file"
	"gin_example/pkg/setting"
	"gin_example/pkg/util"
	"image/jpeg"
	"path"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

type QrCode struct {
	URL    string
	Width  int
	Height int
	Ext    string
	Level  qr.ErrorCorrectionLevel
	Mode   qr.Encoding
}

const (
	EXT_JPG = ".jpg"
)

func NewQrCode(
	url string, with, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {

	return &QrCode{
		URL:    url,
		Width:  with,
		Height: height,
		Ext:    EXT_JPG,
		Level:  level,
		Mode:   mode,
	}
}

func GetQrCodePath() string {
	return setting.AppSetting.QrCodeSavePath
}

func GetQrCodeFullPath() string {
	return path.Join(setting.AppSetting.RuntimeRootPath, setting.AppSetting.QrCodeSavePath)
}

func GetQrCodeFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetQrCodePath() + name
}

// 获得MD5编码后的文件名
func GetQrCodeFileName(value string) string {
	return util.EncodeMD5(value)
}

// 获取文件类型
func (q *QrCode) GetExt() string {
	return q.Ext
}

func (q *QrCode) CheckEncode(path string) bool {
	src := path + GetQrCodeFileName(q.URL) + q.GetExt()
	return !file.CheckNotExist(src)
}

// 编码并保存为二维码
func (q *QrCode) Encode(filePath string) (string, string, error) {
	name := GetQrCodeFileName(q.URL) + q.GetExt()
	src := filePath + name
	if file.CheckNotExist(src) {
		code, err := qr.Encode(q.URL, q.Level, q.Mode)
		if err != nil {
			return "", "", err
		}

		code, err = barcode.Scale(code, q.Width, q.Height)
		if err != nil {
			return "", "", err
		}

		f, err := file.MustOpen(name, filePath)
		if err != nil {
			return "", "", err
		}
		defer f.Close()

		err = jpeg.Encode(f, code, nil)
		if err != nil {
			return "", "", err
		}
	}

	return name, filePath, nil
}
