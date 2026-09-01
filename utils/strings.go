package utils

import "strings"

func NormalizeServiceName(name string) string {
	normalized := strings.ToLower(strings.TrimSpace(name))
	if normalized == "" {
		return "socks5"
	}
	return normalized
}

func NormalizeRuntimeMode(mode string) string {
	return strings.ToLower(strings.TrimSpace(mode))
}

func IsServerMode(mode string) bool {
	normalized := NormalizeRuntimeMode(mode)
	return normalized != "client"
}
