package utils

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	mathrand "math/rand"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用bcrypt加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomPassword 生成随机密码
func GenerateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[mathrand.Intn(len(charset))]
	}
	return string(b)
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	cryptorand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
}
