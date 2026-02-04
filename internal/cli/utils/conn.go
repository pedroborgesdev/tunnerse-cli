package utils

import (
	"errors"
	"net"
	"strings"
)

// IsConnRefused detecta erro de conexão recusada (API local offline)
func IsConnRefused(err error) bool {
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		if netErr.Err != nil && strings.Contains(netErr.Err.Error(), "connection refused") {
			return true
		}
	}
	if strings.Contains(err.Error(), "connection refused") {
		return true
	}
	return false
}
