package app

import (
	"gin_example/pkg/logging"

	"github.com/astaxie/beego/validation"
)

// 统一写出参数校验失败
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}
}
