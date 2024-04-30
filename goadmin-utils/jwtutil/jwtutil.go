package jwtutil

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

//var jwtSignKey = []byte(global.GVA_CONFIG.System.JWT_SIGN_KEY)

//var jwtSignKey = []byte("defalut")

// 1.自定义声明类型
type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 2.封装生成 token 的函数
func GenToken(username string, jwtSignKey string, jwtExpriReTime int64) (string, error) {
	claims := MyCustomClaims{
		username,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(jwtExpriReTime))), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()), // 签名生效时间
			Issuer:    "luoyexinsheng",                // 签名的发行者
			Subject:   "ljzh",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(jwtSignKey))
	return ss, err
}

// 3. 解析 token
func ParseToken(tokenString string, jwtSignKey string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSignKey), nil
	})
	if err != nil {
		// 解析token失败
		//global.GVA_LOG.Error("解析Token失败")
		return nil, err
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		// 说明token合法
		return claims, nil
	} else {
		// token 不合法
		//global.GVA_LOG.Warning("Token不合法")
		return nil, errors.New("Token不合法: invalid token")
	}
}
