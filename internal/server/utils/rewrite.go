package utils

import (
	"fmt"
	"strings"
)


func RewriteAbsolutePaths(html []byte, tunnelName string) []byte {
	content := string(html)
	prefix := fmt.Sprintf("/%s", tunnelName)

	content = strings.ReplaceAll(content, `href="/`, fmt.Sprintf(`href="%s/`, prefix))
	content = strings.ReplaceAll(content, `src="/`, fmt.Sprintf(`src="%s/`, prefix))
	content = strings.ReplaceAll(content, `action="/`, fmt.Sprintf(`action="%s/`, prefix))

	return []byte(content)
}


func InjectBaseHref(body []byte, tunnelName string) []byte {
	html := string(body)
	baseTag := fmt.Sprintf(`<base href="/%s/">`, tunnelName)
	html = strings.Replace(html, "<head>", "<head>\n"+baseTag, 1)
	return []byte(html)
}
