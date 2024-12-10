package util

import (
	"errors"
	"gin_example/pkg/setting"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var TOKEN_COOKIE_KEY string = "_userToken"
var jwtSecret = []byte(setting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

// 签发一个jwt
func GenerateToken(username, password string) (tokenString string, err error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Hour * 3)

	claim := Claims{
		Username: username,
		Password: password,
		RegisteredClaims: jwt.RegisteredClaims{
			// 签发时间
			IssuedAt: jwt.NewNumericDate(nowTime),
			// 过期时间
			ExpiresAt: jwt.NewNumericDate(expireTime),
			// 设置生效时间为当前
			NotBefore: jwt.NewNumericDate(nowTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err = token.SignedString(jwtSecret)
	if err != nil {
		log.Println("token 生成失败", tokenString, err)
	}
	return
}

// JWT TOKEN 解码
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil // 返回签名时使用的密钥
	})

	if err != nil {
		return nil, err
	}
	claim, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("couldn't handle this token")
	}
	return claim, nil
}
