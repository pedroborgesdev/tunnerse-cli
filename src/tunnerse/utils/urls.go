package utils

import "tunnerse/config"

// GetUrl returns the appropriate server URL based on the requested method type.
func GetUrl(method string) string {
	switch method {
	case "register":
		return "https://" + config.GetServerURL() + "/register"
	case "response":
		return config.GetTunnelHTTPURL() + "/response"
	case "fetch":
		return config.GetTunnelHTTPURL() + "/tunnel"
	case "ping":
		return config.GetTunnelHTTPURL()
	}
	return "undefined"
}
