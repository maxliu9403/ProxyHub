package logic

import (
	"net"
	"strings"
)

func CheckIP(ip string) bool {
	if net.ParseIP(ip) == nil || strings.Contains(ip, ":") {
		return false
	}
	return true
}
