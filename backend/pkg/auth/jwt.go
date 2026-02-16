package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired     = errors.New("token已过期")
	ErrTokenInvalid     = errors.New("token无效")
	ErrTokenMalformed   = errors.New("token格式错误")
	ErrTokenNotValidYet = errors.New("token尚未生效")
)

// Claims JWT claims结构
type Claims struct {
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	RealName     string `json:"real_name"`
	Organization string `json:"organization,omitempty"`
	Department   string `json:"department,omitempty"`
	Roles        string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string
	ExpireTime time.Duration
}

// JWT JWT工具结构
type JWT struct {
	config *JWTConfig
}

// NewJWT 创建JWT实例
func NewJWT(config *JWTConfig) *JWT {
	return &JWT{config: config}
}

// GenerateToken 生成token
func (j *JWT) GenerateToken(userID uint, username, realName, organization, department, roles string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:       userID,
		Username:     username,
		RealName:     realName,
		Organization: organization,
		Department:   department,
		Roles:        roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.config.ExpireTime)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "plm-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.Secret))
}

// ParseToken 解析token
func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshToken 刷新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	return j.GenerateToken(
		claims.UserID,
		claims.Username,
		claims.RealName,
		claims.Organization,
		claims.Department,
		claims.Roles,
	)
}

// GetUserIDFromToken 从token中获取用户ID
func (j *JWT) GetUserIDFromToken(tokenString string) (uint, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
