package commands_utils

import (
	"fmt"
	"tunnerse/config"
)

// buildTunnelURL gera a URL de acesso ao túnel com ou sem subdomínio.
func BuildTunnelURL() string {
	if config.GetSubdomainBool() {
		return fmt.Sprintf("https://%s.%s", config.GetTunnelID(), config.GetServerURL())
	}
	return fmt.Sprintf("https://%s/%s", config.GetServerURL(), config.GetTunnelID())
}
