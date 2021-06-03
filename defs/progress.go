package defs

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

type JSONPingProgress struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Ping      struct {
		Jitter   float64 `json:"jitter"`
		Latency  float64 `json:"latency"`
		Progress float64 `json:"progress"`
	} `json:"ping"`
}

type JSONDownloadProgress struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Download  struct {
		Bandwidth int     `json:"bandwidth"`
		Bytes     int     `json:"bytes"`
		Elapsed   int64   `json:"elapsed"`
		Progress  float64 `json:"progress"`
	} `json:"download"`
}

type JSONUploadProgress struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Upload    struct {
		Bandwidth int     `json:"bandwidth"`
		Bytes     int     `json:"bytes"`
		Elapsed   int64   `json:"elapsed"`
		Progress  float64 `json:"progress"`
	} `json:"upload"`
}

func sendPingProgress(latency float64, jitter float64, progress float64) {
	var pingProgress JSONPingProgress
	pingProgress.Timestamp = time.Now()
	pingProgress.Type = "ping"
	pingProgress.Ping.Latency = latency
	pingProgress.Ping.Jitter = jitter
	pingProgress.Ping.Progress = progress

	if b, err := json.Marshal(&pingProgress); err != nil {
		log.Errorf("Error generating progress update: %s", err)
	} else {
		log.Warnf("%s", b)
	}
}

func sendDownloadProgress(c *BytesCounter, durationMs int64) {
	var progress JSONDownloadProgress

	progress.Timestamp = time.Now()
	progress.Type = "download"
	progress.Download.Bytes = c.total
	progress.Download.Elapsed = time.Since(c.start).Milliseconds()
	progress.Download.Bandwidth = int(float64(progress.Download.Bytes) / (float64(time.Since(c.start).Milliseconds()) / 1000))
	progress.Download.Progress = float64(progress.Download.Elapsed) / float64(durationMs)
	if progress.Download.Progress > 1 {
		progress.Download.Progress = 1
	}

	if b, err := json.Marshal(&progress); err != nil {
		log.Errorf("Error generating progress update: %s", err)
	} else {
		log.Warnf("%s", b)
	}
}

// TODO: set durationMs once instead of with every progress update since it's constant
func sendUploadProgress(c *BytesCounter, durationMs int64) {
	var progress JSONUploadProgress

	progress.Timestamp = time.Now()
	progress.Type = "upload"
	progress.Upload.Bytes = c.total
	progress.Upload.Elapsed = time.Since(c.start).Milliseconds()
	progress.Upload.Bandwidth = int(float64(progress.Upload.Bytes) / (float64(time.Since(c.start).Milliseconds()) / 1000))
	progress.Upload.Progress = float64(progress.Upload.Elapsed) / float64(durationMs)
	if progress.Upload.Progress > 1 {
		progress.Upload.Progress = 1
	}

	if b, err := json.Marshal(&progress); err != nil {
		log.Errorf("Error generating progress update: %s", err)
	} else {
		log.Warnf("%s", b)
	}
}
