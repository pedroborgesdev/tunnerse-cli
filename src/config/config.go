package config

import (
	"fmt"
	"sync"
)

var (
	tunnelID   = "undefined"
	serverURL  = "tunnerse.com"
	addressURL = "127.0.0.1:5500"
	mu         sync.RWMutex
)

// SetTunnelID updates the tunnel ID
func SetTunnelID(id string) {
	mu.Lock()
	defer mu.Unlock()
	tunnelID = id
}

// SetServerURL updates the server domain
func SetServerURL(url string) {
	mu.Lock()
	defer mu.Unlock()
	serverURL = url
}

// SetAddressURL updates the address URL based on input
func SetAddressURL(address string) {
	mu.Lock()
	defer mu.Unlock()
	addressURL = fmt.Sprintf("http://127.0.0.1:" + address)
}

// GetTunnelID returns the current tunnel ID
func GetTunnelID() string {
	mu.RLock()
	defer mu.RUnlock()
	return tunnelID
}

// GetServerURL returns the current server domain
func GetServerURL() string {
	mu.RLock()
	defer mu.RUnlock()
	return serverURL
}

// GetTunnelURL returns the complete HTTP tunnel URL
func GetTunnelURL() string {
	mu.RLock()
	defer mu.RUnlock()
	return "http://" + tunnelID + "." + serverURL
}

// GetTunnelHTTPSURL returns the complete HTTPS tunnel URL
func GetTunnelHTTPSURL() string {
	mu.RLock()
	defer mu.RUnlock()
	return "http://" + tunnelID + "." + serverURL
}

// GetAddressURL returns the current address URL
func GetAddressURL() string {
	mu.RLock()
	defer mu.RUnlock()
	return addressURL
}
