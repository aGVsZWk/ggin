package util

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"ggin/pkg/setting"
)

var jwtSecret = []byte(setting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)
	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "ggin",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)	// 按指定方法加密 生成 tokenClaims
	token, err := tokenClaims.SignedString(jwtSecret)	// 对 tokenClaims 加签名，生成 token

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token)(interface{}, error){
	return jwtSecret, nil
	})		// 解析 token 生成 tokenClaims

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {  // 验证 tokenClaims 签名（是否过期）
			return claims, nil
		}
	}
	return nil, err


}
