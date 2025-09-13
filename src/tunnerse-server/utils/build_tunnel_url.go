package utils

import (
	"fmt"
)

// buildTunnelURL gera a URL de acesso ao túnel com ou sem subdomínio.
func BuildTunnelURL(ID, ServerURL string, isSubdomain bool) string {
	if isSubdomain {
		return fmt.Sprintf("https://%s.%s", ID, ServerURL)
	}
	return fmt.Sprintf("https://%s/%s", ServerURL, ID)
}
