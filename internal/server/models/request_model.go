package models

import (
	"encoding/base64"
	"encoding/json"
)

type RequestData struct {
	Method    string              `json:"method"`
	Path      string              `json:"path"`
	Headers   map[string][]string `json:"headers"`
	Body      string              `json:"body"`
	Host      string              `json:"host"`
	RequestID string              `json:"request_id"`
	Token     string              `json:"token"` // Tunnerse-Request-Token
}

type ResponseData struct {
	StatusCode int                 `json:"status_code"`
	Headers    map[string][]string `json:"headers"`
	Body       []byte              `json:"-"`     // Não serializa diretamente
	Token      string              `json:"token"` // Tunnerse-Request-Token
}

// MarshalJSON customiza a serialização para converter Body em base64
func (r *ResponseData) MarshalJSON() ([]byte, error) {
	type Alias struct {
		StatusCode int                 `json:"status_code"`
		Headers    map[string][]string `json:"headers"`
		Body       string              `json:"body"`
		Token      string              `json:"token"`
	}

	return json.Marshal(&Alias{
		StatusCode: r.StatusCode,
		Headers:    r.Headers,
		Body:       base64.StdEncoding.EncodeToString(r.Body),
		Token:      r.Token,
	})
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
