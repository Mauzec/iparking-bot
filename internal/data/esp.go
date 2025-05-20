package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mauzec/ibot-things/pkg/logger"
)

var (
	ESPClientErrNotConnected = fmt.Errorf("unable to connect")
)

type lastDistances [5]float64

func (d *lastDistances) PushAt(idx int, v float64) {
	(*d)[idx] = v
}
func (d *lastDistances) GetAvg() float64 {
	sum, good := float64(0), float64(0)
	for _, v := range *d {
		if v > 0 {
			sum += v
			good++
		}
	}
	if good == 0 {
		return 0.0
	}
	return sum / good
}
func (d *lastDistances) GetMax() float64 {
	max := float64(0)
	for _, v := range *d {
		if v > max {
			max = v
		}
	}
	return max
}
func (d *lastDistances) GetAt(idx int) float64 {
	return (*d)[idx]
}

type EspProvider struct {
	Addr     string
	Port     int
	Interval time.Duration

	mu sync.RWMutex

	lastDistances lastDistances
	lastDistIdx   int // 0..4

	DataPath string
}

func NewEspProvider(addr string, port int, dataPath string, interval time.Duration) *EspProvider {
	if addr == "" || port <= 0 || dataPath == "" {
		logger.Logger().Panicf("[EspProvider] invalid parameters: addr=%s, port=%d, dataPath=%s", addr, port, dataPath)
	}

	return &EspProvider{
		Addr:          addr,
		Port:          port,
		lastDistances: lastDistances{},
		lastDistIdx:   0,
		DataPath:      dataPath,
		Interval:      interval,
	}
}

func (p *EspProvider) Go(ctx context.Context) {
	go p.run(ctx)
}

func (p *EspProvider) run(ctx context.Context) {
	if err := os.MkdirAll(filepath.Dir(p.DataPath), 0o755); err != nil {
		logger.Logger().Panicf("[EspProvider] mkdir: %v", err)
	}

	url := "http://" + p.Addr + ":" + strconv.Itoa(p.Port) + "/"
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	for {
		var err error
		var v float64
		var s string

		var body []byte
		var resp *http.Response

		select {
		case <-ctx.Done():
			logger.Logger().Info("[EspProvider] stopping")
			return
		default:
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			logger.Logger().Warnf("[EspProvider] build request: %v", err)
			goto SLEEP
		}
		req.Close = true

		resp, err = client.Do(req)
		if err != nil {
			logger.Logger().Warnf("[EspProvider] request error: %v", err)
			goto SLEEP
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			if strings.Contains(err.Error(), "reset by peer") {
				// esp library always close socket after request

				// logger.Logger().Debugf("[EspProvider] connection reset by peer, retrying")
			} else {
				logger.Logger().Warnf("[EspProvider] read body: %v", err)
				goto SLEEP
			}

		}
		err = resp.Body.Close()
		if err != nil {
			logger.Logger().Warnf("[EspProvider] close body: %v", err)
		}

		s = strings.TrimSpace(string(body))
		s = strings.Split(s, " ")[0]

		// UNCOMMENT if you need to debug, but it will spam the logs, especially if interval is small
		// logger.Logger().Infof("[EspProvider] response body: %q", s)
		v, err = strconv.ParseFloat(s, 64)
		if err != nil {
			logger.Logger().Warnf("[EspProvider, run] parse float: %v", err)
			goto SLEEP
		}

		p.mu.Lock()
		p.lastDistances.PushAt(p.lastDistIdx, v)
		p.lastDistIdx = (p.lastDistIdx + 1) % len(p.lastDistances)
		p.mu.Unlock()

		if err := p.writeJSON(); err != nil {
			logger.Logger().Warnf("[EspProvider] writeJSON: %v", err)
		}

	SLEEP:
		select {
		case <-ctx.Done():
			logger.Logger().Info("[EspProvider] stopping")
			return
		case <-time.After(p.Interval):
		}
	}

}

func (p *EspProvider) writeJSON() error {
	v := p.GetMaxDistance()
	payload := struct {
		Distance float64 `json:"distance"`
	}{Distance: v}

	buf, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	tmp := p.DataPath + ".tmp"
	if err := os.WriteFile(tmp, buf, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, p.DataPath)
}

func (p *EspProvider) GetMaxDistance() float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.lastDistances.GetMax()
}
