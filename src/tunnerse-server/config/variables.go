package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	tunnelID     = "undefined"
	serverURL    = "tunnerse.com"
	addressURL   = "127.0.0.1:5500"
	tunnelsPath  = "./"
	is_subdomain = false
	mu           sync.RWMutex
)

var LocalRoutes = map[string]string{}

// SetExecPath create and set path of executable
func SetExecPath() {
	execPath, _ := os.Executable()

	execDir := filepath.Dir(execPath)

	tunnelsPath = filepath.Join(execDir, "tunnels")

	os.MkdirAll(tunnelsPath, 0755)
}

// SetSubdomainBool updates the subdomain bool
func SetSubdomainBool(subdomain bool) {
	mu.Lock()
	defer mu.Unlock()
	is_subdomain = subdomain
}

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

// GetExecPath returns the executable path
func GetExecPath() string {
	mu.RLock()
	defer mu.RUnlock()
	return tunnelsPath
}

// GetSubdomainBool returns the mode of api accept
func GetSubdomainBool() bool {
	mu.RLock()
	defer mu.RUnlock()
	return is_subdomain
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
func GetTunnelHTTPURL() string {
	mu.RLock()
	defer mu.RUnlock()
	if GetSubdomainBool() {
		return "https://" + tunnelID + "." + serverURL
	}
	return "https://" + serverURL + "/" + tunnelID
}

// GetTunnelHTTPSURL returns the complete HTTPS tunnel URL
func GetTunnelHTTPSURL() string {
	mu.RLock()
	defer mu.RUnlock()
	if GetSubdomainBool() {
		return "https://" + tunnelID + "." + serverURL
	}
	return "https://" + serverURL + "/" + tunnelID
}

// GetAddressURL returns the current address URL
func GetAddressURL() string {
	mu.RLock()
	defer mu.RUnlock()
	return addressURL
}
