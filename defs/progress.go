package defs

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type JSONProgressInterfaceInfo struct {
	Name       string `json:"name"`
	MacAddr    string `json:"macAddr"`
	IsVpn      bool   `json:"isVpn"`
	ExternalIP string `json:"externalIp"`
	InternalIP string `json:"internalIp"`
}

type JSONProgressServerInfo struct {
	Name     string `json:"name"`
	Country  string `json:"country"`
	Host     string `json:"host"`
	Location string `json:"location"`
	ID       int    `json:"id"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
}

type JSONProgressHeader struct {
	Type      string                    `json:"type"`
	Timestamp time.Time                 `json:"timestamp"`
	ISP       string                    `json:"isp"`
	Server    JSONProgressServerInfo    `json:"server"`
	Interface JSONProgressInterfaceInfo `json:"interface"`
}

type JSONProgressPing struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Ping      struct {
		Jitter   float64 `json:"jitter"`
		Latency  float64 `json:"latency"`
		Progress float64 `json:"progress"`
	} `json:"ping"`
}

type JSONProgressDownload struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Download  struct {
		Bandwidth int     `json:"bandwidth"`
		Bytes     int     `json:"bytes"`
		Elapsed   int64   `json:"elapsed"`
		Progress  float64 `json:"progress"`
	} `json:"download"`
}

type JSONProgressUpload struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Upload    struct {
		Bandwidth int     `json:"bandwidth"`
		Bytes     int     `json:"bytes"`
		Elapsed   int64   `json:"elapsed"`
		Progress  float64 `json:"progress"`
	} `json:"upload"`
}

type InterfaceStats struct {
	RxBytes      int64 `json:"rxbytes"`
	TxBytes      int64 `json:"txbytes"`
	RxPackets    int64 `json:"rxpackets"`
	TxPackets    int64 `json:"txpackets"`
	RxErrors     int64 `json:"rxerrors"`
	TxErrors     int64 `json:"txerrors"`
	RxDropped    int64 `json:"rxdropped"`
	TxDropped    int64 `json:"txdropped"`
	RxFifo       int64 `json:"rxfifo"`
	TxFifo       int64 `json:"txfifo"`
	RxFrame      int64 `json:"rxframe"`
	TxFrame      int64 `json:"txframe"`
	RxCompressed int64 `json:"rxcompressed"`
	TxCompressed int64 `json:"txcompressed"`
	RxMulticast  int64 `json:"rxmulticast"`
	TxMulticast  int64 `json:"txmulticast"`
}

func getInterfaceStats(i *net.Interface) InterfaceStats {
	var result InterfaceStats
	filename := "/Users/chris/src/cpec/speedtest-cli/ifstat.txt"

	file, err := os.Open(filename)
	if err != nil {
		log.Errorf("error: Failed to open %s (%s)", filename, err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "eth0:") {
			fields := strings.Fields(line)
			rxfields := fields[1:9]
			if i, err := strconv.ParseInt(rxfields[0], 10, 64); err == nil {
				result.RxBytes = (i)
			}
			if i, err := strconv.ParseInt(rxfields[1], 10, 64); err == nil {
				result.RxPackets = (i)
			}
			if i, err := strconv.ParseInt(rxfields[2], 10, 64); err == nil {
				result.RxErrors = (i)
			}
			if i, err := strconv.ParseInt(rxfields[3], 10, 64); err == nil {
				result.RxDropped = (i)
			}
			if i, err := strconv.ParseInt(rxfields[4], 10, 64); err == nil {
				result.RxFifo = (i)
			}
			if i, err := strconv.ParseInt(rxfields[5], 10, 64); err == nil {
				result.RxFrame = (i)
			}
			if i, err := strconv.ParseInt(rxfields[6], 10, 64); err == nil {
				result.RxCompressed = (i)
			}
			if i, err := strconv.ParseInt(rxfields[7], 10, 64); err == nil {
				result.RxMulticast = (i)
			}
			txfields := fields[9:]
			if i, err := strconv.ParseInt(txfields[0], 10, 64); err == nil {
				result.TxBytes = (i)
			}
			if i, err := strconv.ParseInt(txfields[1], 10, 64); err == nil {
				result.TxPackets = (i)
			}
			if i, err := strconv.ParseInt(txfields[2], 10, 64); err == nil {
				result.TxErrors = (i)
			}
			if i, err := strconv.ParseInt(txfields[3], 10, 64); err == nil {
				result.TxDropped = (i)
			}
			if i, err := strconv.ParseInt(txfields[4], 10, 64); err == nil {
				result.TxFifo = (i)
			}
			if i, err := strconv.ParseInt(txfields[5], 10, 64); err == nil {
				result.TxFrame = (i)
			}
			if i, err := strconv.ParseInt(txfields[6], 10, 64); err == nil {
				result.TxCompressed = (i)
			}
			if i, err := strconv.ParseInt(txfields[7], 10, 64); err == nil {
				result.TxMulticast = (i)
			}
		}
	}

	return result
}

func SendProgressHeader(s *Server, isp *IPInfoResponse) {
	getWanInterface := func() net.Interface {
		var result net.Interface

		ifs, _ := net.Interfaces()
		for i := 0; i < len(ifs); i++ {
			if ifs[i].Name == "eth0" {
				result = ifs[i]
				break
			} else if ifs[i].Flags&net.FlagUp != 0 && ifs[i].Flags&net.FlagLoopback == 0 {
				ips, _ := ifs[i].Addrs()
				if len(ips) > 1 {
					result = ifs[i]
					break
				}
			}
		}

		return result
	}

	getIPFromInterface := func(iface *net.Interface) string {
		// find the wan IP address
		var ipAddr net.IP
		if addrs, err := iface.Addrs(); err == nil {
			for _, addr := range addrs {
				if ipAddr = addr.(*net.IPNet).IP.To4(); ipAddr != nil {
					return ipAddr.String()
				}
			}
		}

		return "n/a"
	}

	var header JSONProgressHeader
	wanInterface := getWanInterface()
	stats := getInterfaceStats(&wanInterface)
	if b, err := json.Marshal(&stats); err != nil {
		log.Errorf("Error serializing interface stats: %s", err)
	} else {
		log.Warnf("%s", b)
	}

	header.Timestamp = time.Now()
	header.Type = "testStart"
	header.ISP = isp.Organization
	header.Interface.ExternalIP = isp.IP
	header.Interface.InternalIP = getIPFromInterface(&wanInterface)
	header.Interface.IsVpn = false
	header.Interface.MacAddr = wanInterface.HardwareAddr.String()
	header.Interface.Name = wanInterface.Name

	serverUrl, _ := s.GetURL()
	header.Server.Name = s.Name
	header.Server.ID = s.ID
	header.Server.Host = serverUrl.Hostname()
	header.Server.Port = serverUrl.Port()
	header.Server.IP = serverUrl.Hostname()
	header.Server.Country = s.Country
	header.Server.Location = s.Location

	if b, err := json.Marshal(&header); err != nil {
		log.Errorf("Error generating progress update: %s", err)
	} else {
		log.Warnf("%s", b)
	}
}

func SendPingProgress(latency float64, jitter float64, progress float64) {
	var pingProgress JSONProgressPing
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

func SendDownloadProgress(c *BytesCounter, durationMs int64) {
	var progress JSONProgressDownload

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
func SendUploadProgress(c *BytesCounter, durationMs int64) {
	var progress JSONProgressUpload

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
