package utils

import (
	"fmt"
	"strings"
)

type RewriteUtils struct{}

// NewRewriteUtils creates and returns a new instance of RewriteUtils.
func NewRewriteUtils() *RewriteUtils {
	return &RewriteUtils{}
}

// RewriteAbsolutePaths rewrites all href, src, and action attributes to include the tunnel name as a prefix.
func (u *RewriteUtils) RewriteAbsolutePaths(html []byte, tunnelName string) []byte {
	content := string(html)
	prefix := fmt.Sprintf("/%s", tunnelName)

	content = strings.ReplaceAll(content, `href="/`, fmt.Sprintf(`href="%s/`, prefix))
	content = strings.ReplaceAll(content, `src="/`, fmt.Sprintf(`src="%s/`, prefix))
	content = strings.ReplaceAll(content, `action="/`, fmt.Sprintf(`action="%s/`, prefix))

	return []byte(content)
}

// InjectBaseHref injects a <base> tag into the <head> to set the base URL with the tunnel name.
func (u *RewriteUtils) InjectBaseHref(body []byte, tunnelName string) []byte {
	html := string(body)
	baseTag := fmt.Sprintf(`<base href="/%s/">`, tunnelName)
	html = strings.Replace(html, "<head>", "<head>\n"+baseTag, 1)
	return []byte(html)
}
