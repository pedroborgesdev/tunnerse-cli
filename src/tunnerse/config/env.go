package config

import "os"

func GetEnvBackgroundRunning() bool {
	if os.Getenv("TUNNERSE_BG") == "1" {
		return true
	} else {
		return false
	}
}

func GetEnvApplicationRunning() bool {
	if os.Getenv("TUNNERSE_APP") == "1" {
		return true
	} else {
		return false
	}
}

func GetEnvTunneID() string {
	return os.Getenv("TUNNERSE_TUNNEL_ID")
}

func GetEnvSubdomain() string {
	return os.Getenv("TUNNERSE_SUBDOMAIN")
}

func GetEnvServerUrl() string {
	return os.Getenv("TUNNERSE_SERVER_URL")
}

func GetEnvAddressPort() string {
	return os.Getenv("TUNNERSE_ADDRESS_PORT")
}
