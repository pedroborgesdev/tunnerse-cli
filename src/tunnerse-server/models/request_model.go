package models

// RequestData represents an HTTP request received via the tunnel, including method, path, headers, and body.
type RequestData struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
	Host    string              `json:"host"`
}

// ResponseData represents the HTTP response to be sent back through the tunnel, including status, headers, and body.
type ResponseData struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

// RegisterResponse represents the response structure received after registering a tunnel with the server.
type RegisterResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Message   string `json:"message"`
		Subdomain bool   `json:"subdomain"`
		Tunnel    string `json:"tunnel"`
	} `json:"data"`
	Status int `json:"status"`
}
