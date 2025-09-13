package utils

type RegisterRequest struct {
	Name string `json:"name" binding:"required"`
}

type OpenRequest struct {
	Name      string `json:"name" binding:"required"`
	Port      string `json:"port" binding:"required"`
	ServerURL string `json:"server_domain" binding:"required"`
}
