package middleware

import (
	"gin_example/pkg/e"
	log "gin_example/pkg/logging"
	"gin_example/pkg/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 校验jwt是否存在
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		// tokenStr := c.Query("token")
		token, err := c.Request.Cookie(util.TOKEN_COOKIE_KEY)
		var tokenStr string = ""
		if err == nil {
			tokenStr = token.Value
		}
		// 对Cookie有效期进行验证
		if tokenStr == "" {
			code = e.ERROR_NOTFOUND_TOKEN
		} else {
			// 解析并坚定
			claim, err := util.ParseToken(tokenStr)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claim.ExpiresAt.Unix() {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		// 失败鉴权处理
		if code != e.SUCCESS {
			log.InfoF("refuse %s request with Token [%s], message:[%s]",
				c.Request.RemoteAddr, tokenStr, e.GetMsg(code))
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})
			c.Abort() // 终止后续请求
			return
		}
		// 成功鉴权可以走后续的处理逻辑
		// 其实不用调用Next是不是更好，因为后续就可以自己按照前面的流程走下去了，可以少走一个循环
		// c.Next()
	}
}
