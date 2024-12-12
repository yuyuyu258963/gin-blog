package middleware

import (
	log "gin_example/pkg/logging"
	"time"

	"github.com/gin-gonic/gin"
)

// 日志中间件
func Logger() gin.HandlerFunc {

	return func(c *gin.Context) {
		startTime := time.Now()
		// 处理后续请求
		c.Next()
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求Ip
		ClientIP := c.ClientIP()
		log.InfoFiled(log.LogFields{
			"latencyTime": latencyTime,
			"reqMethod":   reqMethod,
			"reqUri":      reqUri,
			"statusCode":  statusCode,
			"ClientIP":    ClientIP,
		})
	}
}
