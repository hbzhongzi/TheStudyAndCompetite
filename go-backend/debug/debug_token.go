package main

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 复制Claims结构体以避免导入问题
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("yunmeng-secret-key")

// 生成JWT Token
func GenerateToken(userID uint, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7天有效期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()), // 立即生效
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// 解析JWT Token
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

func main() {
	log.Println("=== JWT Token调试工具 ===")

	// 测试Token生成
	log.Println("1. 测试Token生成...")
	testTokenGeneration()

	// 测试Token解析
	log.Println("\n2. 测试Token解析...")
	testTokenParsing()

	// 测试Token过期时间
	log.Println("\n3. 测试Token过期时间...")
	testTokenExpiration()

	// 测试Token验证
	log.Println("\n4. 测试Token验证...")
	testTokenValidation()
}

func testTokenGeneration() {
	// 生成测试Token
	token, err := GenerateToken(1, "admin", "admin")
	if err != nil {
		log.Printf("❌ Token生成失败: %v", err)
		return
	}

	log.Printf("✅ Token生成成功")
	log.Printf("Token长度: %d", len(token))
	log.Printf("Token前50字符: %s...", token[:50])

	// 手动解析Token查看内容
	parts, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("yunmeng-secret-key"), nil
	})

	if err != nil {
		log.Printf("❌ Token解析失败: %v", err)
		return
	}

	if claims, ok := parts.Claims.(*Claims); ok {
		log.Printf("✅ Token解析成功")
		log.Printf("用户ID: %d", claims.UserID)
		log.Printf("用户名: %s", claims.Username)
		log.Printf("角色: %s", claims.Role)
		log.Printf("签发时间: %v", claims.IssuedAt.Time)
		log.Printf("过期时间: %v", claims.ExpiresAt.Time)
		log.Printf("当前时间: %v", time.Now())
		log.Printf("是否已过期: %v", time.Now().After(claims.ExpiresAt.Time))
	}
}

func testTokenParsing() {
	// 生成一个Token
	token, err := GenerateToken(1, "admin", "admin")
	if err != nil {
		log.Printf("❌ Token生成失败: %v", err)
		return
	}

	// 使用ParseToken解析
	claims, err := ParseToken(token)
	if err != nil {
		log.Printf("❌ Token解析失败: %v", err)
		return
	}

	log.Printf("✅ Token解析成功")
	log.Printf("用户ID: %d", claims.UserID)
	log.Printf("用户名: %s", claims.Username)
	log.Printf("角色: %s", claims.Role)
	log.Printf("签发时间: %v", claims.IssuedAt.Time)
	log.Printf("过期时间: %v", claims.ExpiresAt.Time)
}

func testTokenExpiration() {
	// 测试不同过期时间的Token
	testCases := []struct {
		name     string
		duration time.Duration
	}{
		{"1分钟", 1 * time.Minute},
		{"1小时", 1 * time.Hour},
		{"24小时", 24 * time.Hour},
		{"7天", 7 * 24 * time.Hour},
	}

	for _, tc := range testCases {
		log.Printf("\n测试 %s 过期时间:", tc.name)

		// 创建自定义过期时间的Token
		claims := Claims{
			UserID:   1,
			Username: "admin",
			Role:     "admin",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(tc.duration)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("yunmeng-secret-key"))
		if err != nil {
			log.Printf("❌ Token生成失败: %v", err)
			continue
		}

		// 解析Token
		parsedClaims, err := ParseToken(tokenString)
		if err != nil {
			log.Printf("❌ Token解析失败: %v", err)
			continue
		}

		log.Printf("✅ Token生成和解析成功")
		log.Printf("  过期时间: %v", parsedClaims.ExpiresAt.Time)
		log.Printf("  当前时间: %v", time.Now())
		log.Printf("  是否已过期: %v", time.Now().After(parsedClaims.ExpiresAt.Time))
		log.Printf("  剩余时间: %v", parsedClaims.ExpiresAt.Time.Sub(time.Now()))
	}
}

// 测试Token验证中间件
func testTokenValidation() {
	// 生成有效Token
	validToken, _ := GenerateToken(1, "admin", "admin")

	// 生成过期Token（1秒后过期）
	expiredClaims := Claims{
		UserID:   1,
		Username: "admin",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, _ := expiredToken.SignedString([]byte("yunmeng-secret-key"))

	// 等待Token过期
	time.Sleep(2 * time.Second)

	// 测试有效Token
	_, err := ParseToken(validToken)
	if err != nil {
		log.Printf("❌ 有效Token验证失败: %v", err)
	} else {
		log.Printf("✅ 有效Token验证成功")
	}

	// 测试过期Token
	_, err = ParseToken(expiredTokenString)
	if err != nil {
		log.Printf("✅ 过期Token正确被拒绝: %v", err)
	} else {
		log.Printf("❌ 过期Token未被拒绝")
	}
}
