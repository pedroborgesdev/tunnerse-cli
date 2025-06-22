package utils

import "tunnerse/config"

type UrlsUtils struct{}

// NewUrlsUtils creates and returns a new instance of UrlsUtils.
func NewUrlsUtils() *UrlsUtils {
	return &UrlsUtils{}
}

// GetUrl returns the appropriate server URL based on the requested method type.
func (s *UrlsUtils) GetUrl(method string) string {
	switch method {
	case "register":
		return "http://" + config.GetServerURL() + "/register"
	case "response":
		return config.GetTunnelHTTPSURL() + "/response"
	case "fetch":
		return config.GetTunnelHTTPSURL() + "/tunnel"
	case "ping":
		return config.GetTunnelHTTPSURL()
	}
	return "undefined"
}
