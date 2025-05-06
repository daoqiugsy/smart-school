package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"smart-school/pkg/config"
	"time"
)

// JWT密钥和过期时间
var (
	jwtSecret   []byte
	jwtDuration time.Duration
)

// InitJWT 初始化JWT配置
func InitJWT(cfg *config.JWTConfig) {
	jwtSecret = []byte(cfg.Secret)
	jwtDuration = cfg.Expire
}

// Claims 自定义JWT声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	UserType int    `json:"user_type"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username string, userType int) (string, error) {
	now := time.Now()
	expirationTime := now.Add(jwtDuration)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "smart-school",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, err
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
