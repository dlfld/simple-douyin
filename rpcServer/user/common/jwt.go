package common

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/douyin/models"
	"time"
)

var jwtKey = []byte("a_secret_crect")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken(user models.User) (string, error) {
	//token有效时间7天
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			//token发放时间
			IssuedAt: time.Now().Unix(),
			Issuer:   "learn.tech",
			Subject:  "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "token发生错误", err
	}
	return tokenString, nil
}

//token使用的加密协议。荷载，储存claims中的信息。前面两部分+jwtKey生成的一个哈希值
//token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjMsImV4cCI6MTY4OTg1ODY3NywiaWF0IjoxNjg5MjUzODc3LCJpc3MiOiJsZWFybi50ZWNoIiwic3ViIjoidXNlciB0b2tlbiJ9.iUVSATejwhiVzW1z10Kc77pvNmff7Mo3_IMYC3LhNhA
//echo eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9 | base64 -D 解密

// ParseToken 解析tokenString
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}
