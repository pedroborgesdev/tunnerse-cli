package models

type RequestData struct {
	Method    string              `json:"method"`
	Path      string              `json:"path"`
	Headers   map[string][]string `json:"headers"`
	Body      string              `json:"body"`
	Host      string              `json:"host"`
	RequestID string              `json:"request_id"`
}

type ResponseData struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

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
