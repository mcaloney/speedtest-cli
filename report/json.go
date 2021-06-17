package report

import (
	"time"

	"librespeed-cli/defs"
)

// JSONReport represents the output data fields in a JSON file
type JSONReport struct {
	Timestamp time.Time                    `json:"timestamp"`
	Server    Server                       `json:"server"`
	Client    Client                       `json:"client"`
	Ping      float64                      `json:"ping"`
	Jitter    float64                      `json:"jitter"`
	Upload    defs.TransferSummaryResponse `json:"upload"`
	Download  defs.TransferSummaryResponse `json:"download"`
	Share     string                       `json:"share"`
}

// Server represents the speed test server's information
type Server struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Location string `json:"location"`
	Country  string `json:"country"`
}

// Client represents the speed test client's information
type Client struct {
	defs.IPInfoResponse
}
