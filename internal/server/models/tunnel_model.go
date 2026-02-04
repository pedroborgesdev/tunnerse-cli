package models

type Tunnel struct {
	ID        string
	Port      string
	Url       string
	Domain    string
	Active    bool
	CreatedAt string
}

type Info struct {
	ID           string
	Requests     int
	Healthchecks int
	Warns        int
	Errors       int
}
