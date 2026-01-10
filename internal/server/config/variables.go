package config

import (
	"sync"
)

var (
	tunnelID     = "undefined"
	serverURL    = "localhost:3000"
	is_subdomain = false
	mu           sync.RWMutex
)

var QuickTunnelURLs = map[string]string{}

type TunnelJob interface {
	Stop()
}

var ActiveJobs = map[string]TunnelJob{}
var jobsMu sync.RWMutex

func SetActiveJob(tunnelID string, job TunnelJob) {
	jobsMu.Lock()
	defer jobsMu.Unlock()
	ActiveJobs[tunnelID] = job
}

func GetActiveJob(tunnelID string) (TunnelJob, bool) {
	jobsMu.RLock()
	defer jobsMu.RUnlock()
	job, exists := ActiveJobs[tunnelID]
	return job, exists
}

func RemoveActiveJob(tunnelID string) {
	jobsMu.Lock()
	defer jobsMu.Unlock()
	delete(ActiveJobs, tunnelID)
}

func SetSubdomainBool(subdomain bool) {
	mu.Lock()
	defer mu.Unlock()
	is_subdomain = subdomain
}

func SetTunnelID(id string) {
	mu.Lock()
	defer mu.Unlock()
	tunnelID = id
}

func SetServerURL(url string) {
	mu.Lock()
	defer mu.Unlock()
	serverURL = url
}

func GetSubdomainBool() bool {
	mu.RLock()
	defer mu.RUnlock()
	return is_subdomain
}

func GetTunnelID() string {
	mu.RLock()
	defer mu.RUnlock()
	return tunnelID
}

func GetServerURL() string {
	mu.RLock()
	defer mu.RUnlock()
	return serverURL
}

func GetTunnelHTTPSURL() string {
	mu.RLock()
	defer mu.RUnlock()
	if GetSubdomainBool() {
		return "https://" + tunnelID + "." + serverURL
	}
	return "https://" + serverURL + "/" + tunnelID
}
