package logic

import (
	"crypto/rand"
	"encoding/base64"
	"net"
	"strings"
)

func CheckIP(ip string) bool {
	if net.ParseIP(ip) == nil || strings.Contains(ip, ":") {
		return false
	}
	return true
}

// GenerateSecureToken returns a secure random string of given byte length.
// 例如传入 32，返回 base64 后长度约为 44 个字符
func GenerateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b), nil
}
