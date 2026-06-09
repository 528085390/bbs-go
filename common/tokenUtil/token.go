package tokenUtil

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateAccessToken 生成 JWT 访问令牌，包含用户 ID 和角色信息，过期时间为 24 小时
func GenerateAccessToken(secret string, userId int64, role []string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"userId": userId,
		"role":   role,
		"iat":    now.Unix(),
		"exp":    now.Add(24 * time.Hour).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(secret))
}

// ValidateAndParseToken 解析并验证JWT token
func ValidateAndParseToken(secret string, tokenString string) (bool, int64, []string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return false, 0, nil, err
	}

	// 验证token是否有效
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// 检查过期时间
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return false, 0, nil, errors.New("token已过期")
			}
		}
	}

	userId, ok := claims["userId"].(float64)
	if !ok {
		return false, 0, nil, errors.New("无法获取用户ID")
	}

	roleInterface, ok := claims["role"].([]interface{})
	if !ok {
		return false, 0, nil, errors.New("无法获取角色信息")
	}

	var roles []string
	for _, role := range roleInterface {
		if roleStr, ok := role.(string); ok {
			roles = append(roles, roleStr)
		}
	}

	return true, int64(userId), roles, nil
}
