package v1

import (
	"gin_example/models"
	"gin_example/pkg/e"
	"gin_example/pkg/util"
	"log"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// @Summary Get Auth Token
// @Description Get Auth Token
// @Param username query string true "Name"
// @Param password query string true "Password"
// @Tags auth
// @Accept json
// @Produce json
// @Success 200  {string} string "ok"
// @Failure 20003 {string} string "ok"
// @Router /api/auth [get]
func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	var tokenStr string
	var err error
	valid := validation.Validation{}
	a := auth{username, password}
	ok, _ := valid.Valid(a)

	code := e.INVALID_PARAMS
	if ok {
		isExist := models.CheckAuth(username, password)
		if isExist {
			tokenStr, err = util.GenerateToken(username, password)

			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				code = e.SUCCESS
			}
		} else {
			code = e.ERROR_AUTH
		}
	} else {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
	}

	// 成功验证登录后设置Cookie
	if code == e.SUCCESS {
		cookie := &http.Cookie{
			Name:     util.TOKEN_COOKIE_KEY,
			Value:    tokenStr,
			Path:     "/",   // 设置Cookie有效访问路径
			HttpOnly: true,  // 防止javascript访问
			Secure:   false, // 现在方便测试测则可以不仅限HTTPS传输
		}
		http.SetCookie(c.Writer, cookie)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}
