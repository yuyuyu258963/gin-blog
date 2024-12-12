package util

import (
	"gin_example/pkg/setting"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// 分页页码的获取方法
func GetPage(c *gin.Context) (result int) {
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * setting.AppSetting.PageSize
	}

	return result
}
