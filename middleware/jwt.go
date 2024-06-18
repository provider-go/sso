package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/provider-go/pkg/logger"
	"github.com/provider-go/pkg/util"
	"time"
)

type InstanceJWT struct {
	SecretKey []byte
}

func InitJwt(SecretKey string) InstanceJWT {
	return InstanceJWT{[]byte(SecretKey)}
}

// CreateClaims
// ID			jwt唯一标识（比如：UUID）
// Issuer		Token 颁发者的唯一标识
// Subject		主题（比如：用户id 或 用户名）
// IssuedAt		签发时间（比如：时间戳）
// Audience		JWT的接收者(应用或平台)
// ExpiresAt	有效期，结束时间 （比如：时间戳）
// NotBefore	有效期，开始时间（比如：时间戳）
func CreateClaims(username string) jwt.RegisteredClaims {
	uuid := util.GetRandString(32)
	return jwt.RegisteredClaims{
		ID:        uuid,
		Issuer:    "Star",
		Subject:   username,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Audience:  jwt.ClaimStrings{"All platforms"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),
	}
}

// GenerateToken 生成token
func (i InstanceJWT) GenerateToken(username string) string {
	// 设置一些声明
	claims := CreateClaims(username)
	// 创建一个新的JWT token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 设置签名并获取token字符串
	token, err := jwtToken.SignedString(i.SecretKey)
	if err != nil {
		return ""
	}
	return token
}

// CreateTokenByOldToken 使用旧token换新token
func (i InstanceJWT) CreateTokenByOldToken(tokenString string) string {
	claims := i.ParseToken(tokenString)
	username, err := claims.GetSubject()
	if err != nil {
		logger.Error("CreateTokenByOldToken", "step", "GetSubject", "err", err)
		return ""
	}
	newClaims := CreateClaims(username)
	// 创建一个新的JWT token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)

	// 设置签名并获取token字符串
	token, err := jwtToken.SignedString(i.SecretKey)
	if err != nil {
		return ""
	}
	return token
}

// ParseToken 解析token
func (i InstanceJWT) ParseToken(tokenString string) jwt.MapClaims {
	// 解析token字符串
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return i.SecretKey, nil
	})

	if err != nil {
		return nil
	}

	// 验证token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims
	}

	return nil
}
