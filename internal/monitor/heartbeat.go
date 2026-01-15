package monitor

import (
	"net/http"
	"time"
)

// Result represents the outcome of a single health check
type Result struct {
	URL        string
	StatusCode int
	Latency    time.Duration
	Timestamp  time.Time
}

// CheckURL performs the actual HTTP request
func CheckURL(url string) Result {
	start := time.Now()
	
	// Set a strict timeout so one slow site doesn't hang our system
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	
	latency := time.Since(start)
	
	if err != nil {
		return Result{URL: url, StatusCode: 0, Latency: latency, Timestamp: time.Now()}
	}
	defer resp.Body.Close()

	return Result{
		URL:        url,
		StatusCode: resp.StatusCode,
		Latency:    latency,
		Timestamp:  time.Now(),
	}
}