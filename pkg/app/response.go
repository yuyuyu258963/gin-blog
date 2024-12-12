package app

import (
	"gin_example/pkg/e"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

// 统一格式响应结果
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"code": errCode,
		"msg":  e.GetMsg(errCode),
		"data": data,
	})
}
