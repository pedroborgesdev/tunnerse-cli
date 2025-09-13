package utils

// GetUrl returns the appropriate server URL based on the requested method type.
func GetUrl(method, ID string, isSubdomain bool, serverDomain string) string {

	var url string
	if method == "fetch" || method == "response" || method == "ping" {
		if isSubdomain {
			url = "https://" + ID + "." + serverDomain
		} else {
			url = "https://" + serverDomain + "/" + ID
		}
	} else if method == "register" {
		url = serverDomain
	}

	switch method {
	case "register":
		return "https://" + url + "/register"
	case "response":
		return url + "/response"
	case "fetch":
		return url + "/tunnel"
	case "ping":
		return url
	}
	return "undefined"
}
