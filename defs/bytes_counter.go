package defs

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"

	//"log"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type JSONProgress struct {
	Bandwidth int     `json:"bandwidth"`
	Bytes     int     `json:"bytes"`
	Elapsed   int64   `json:"elapsed"`
	Progress  float64 `json:"progress"`
}

type JSONDownloadProgress struct {
	Type      string       `json:"type"`
	Timestamp time.Time    `json:"timestamp"`
	Download  JSONProgress `json:"download"`
}

type JSONUploadProgress struct {
	Type      string       `json:"type"`
	Timestamp time.Time    `json:"timestamp"`
	Upload    JSONProgress `json:"upload"`
}

// BytesCounter implements io.Reader and io.Writer interface, for counting bytes being read/written in HTTP requests
type BytesCounter struct {
	start              time.Time
	pos                int
	total              int
	payload            []byte
	reader             io.ReadSeeker
	mebi               bool
	uploadSize         int
	lastProgress       time.Time
	progressIntervalMs int64
	transferType       string
	duration           int64

	lock *sync.Mutex
}

func NewCounter() *BytesCounter {
	return &BytesCounter{
		progressIntervalMs: 100,
		lock:               &sync.Mutex{},
	}
}

// Write implements io.Writer
func (c *BytesCounter) Write(p []byte) (int, error) {
	n := len(p)
	c.lock.Lock()
	c.total += n
	c.lock.Unlock()

	if time.Since(c.lastProgress).Milliseconds() > c.progressIntervalMs && c.transferType == "download" {
		var progress JSONDownloadProgress
		progress.Timestamp = time.Now()
		progress.Type = "download"
		progress.Download.Bytes = c.total
		progress.Download.Elapsed = time.Since(c.start).Milliseconds()
		progress.Download.Bandwidth = int(float64(progress.Download.Bytes) / (float64(time.Since(c.start).Milliseconds()) / 1000))
		progress.Download.Progress = float64(progress.Download.Elapsed) / float64(c.duration)
		if progress.Download.Progress > 1 {
			progress.Download.Progress = 1
		}

		if b, err := json.Marshal(&progress); err != nil {
			log.Errorf("Error generating progress update: %s", err)
		} else {
			log.Warnf("%s", b)
		}
		c.lastProgress = time.Now()
	}

	return n, nil
}

// Read implements io.Reader
func (c *BytesCounter) Read(p []byte) (int, error) {
	n, err := c.reader.Read(p)
	c.lock.Lock()
	c.total += n
	c.pos += n
	if c.pos == c.uploadSize {
		c.resetReader()
	}
	c.lock.Unlock()

	if time.Since(c.lastProgress).Milliseconds() > c.progressIntervalMs && c.transferType == "upload" {
		var progress JSONUploadProgress
		progress.Timestamp = time.Now()
		progress.Type = "upload"
		progress.Upload.Bytes = c.total
		progress.Upload.Elapsed = time.Since(c.start).Milliseconds()
		progress.Upload.Bandwidth = int(float64(progress.Upload.Bytes) / (float64(time.Since(c.start).Milliseconds()) / 1000))
		progress.Upload.Progress = float64(progress.Upload.Elapsed) / float64(c.duration)
		if progress.Upload.Progress > 1 {
			progress.Upload.Progress = 1
		}

		if b, err := json.Marshal(&progress); err != nil {
			log.Errorf("Error generating progress update: %s", err)
		} else {
			log.Warnf("%s", b)
		}
		c.lastProgress = time.Now()
	}

	return n, err
}

// SetBase sets the base for dividing bytes into megabyte or mebibyte
func (c *BytesCounter) SetMebi(mebi bool) {
	c.mebi = mebi
}

// Sets the type of transfer (ping/download/upload)
func (c *BytesCounter) SetTransferType(transferType string) {
	c.transferType = transferType
}

// Set the duration of the test for progress tracking
func (c *BytesCounter) SetDuration(duration int64) {
	c.duration = duration
}

// SetUploadSize sets the size of payload being uploaded
func (c *BytesCounter) SetUploadSize(uploadSize int) {
	c.uploadSize = uploadSize * 1024
}

// AvgBytes returns the average bytes/second
func (c *BytesCounter) AvgBytes() float64 {
	return float64(c.total) / time.Now().Sub(c.start).Seconds()
}

// AvgMbps returns the average mbits/second
func (c *BytesCounter) AvgMbps() float64 {
	var base float64 = 125000
	if c.mebi {
		base = 131072
	}
	return c.AvgBytes() / base
}

// AvgHumanize returns the average bytes/kilobytes/megabytes/gigabytes (or bytes/kibibytes/mebibytes/gibibytes) per second
func (c *BytesCounter) AvgHumanize() string {
	val := c.AvgBytes()

	var base float64 = 1000
	if c.mebi {
		base = 1024
	}

	if val < base {
		return fmt.Sprintf("%.2f bytes/s", val)
	} else if val/base < base {
		return fmt.Sprintf("%.2f KB/s", val/base)
	} else if val/base/base < base {
		return fmt.Sprintf("%.2f MB/s", val/base/base)
	} else {
		return fmt.Sprintf("%.2f GB/s", val/base/base/base)
	}
}

// GenerateBlob generates a random byte array of `uploadSize` in the `payload` field, and sets the `reader` field to
// read from it
func (c *BytesCounter) GenerateBlob() {
	c.payload = getRandomData(c.uploadSize)
	c.reader = bytes.NewReader(c.payload)
}

// resetReader resets the `reader` field to 0 position
func (c *BytesCounter) resetReader() (int64, error) {
	c.pos = 0
	return c.reader.Seek(0, 0)
}

// Start will set the `start` field to current time
func (c *BytesCounter) Start() {
	c.start = time.Now()
	c.lastProgress = c.start
}

// Total returns the total bytes read/written
func (c *BytesCounter) Total() int {
	return c.total
}

// CurrentSpeed returns the current bytes/second
func (c *BytesCounter) CurrentSpeed() float64 {
	return float64(c.total) / time.Now().Sub(c.start).Seconds()
}

// SeekWrapper is a wrapper around io.Reader to give it a noop io.Seeker interface
type SeekWrapper struct {
	io.Reader
}

// Seek implements the io.Seeker interface
func (r *SeekWrapper) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

// getAvg returns the average value of an float64 array
func getAvg(vals []float64) float64 {
	var total float64
	for _, v := range vals {
		total += v
	}

	return total / float64(len(vals))
}

// getRandomData returns an `length` sized array of random bytes
func getRandomData(length int) []byte {
	data := make([]byte, length)
	if _, err := rand.Read(data); err != nil {
		log.Fatalf("Failed to generate random data: %s", err)
	}
	return data
}
