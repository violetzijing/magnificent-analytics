package services

import (
	"magnificent-analytics/lib"
	"time"
)

type MagnificentService struct {
	URL        string
	HTTPClient lib.HTTPClientInterface
	RespQueue  []*CheckResult

	Threshold    int
	HealthRatio  float64
	Timeout      int
	IntervalTime int
}

func (m *MagnificentService) GetIntervalTime() int {
	return m.IntervalTime
}

func NewMagnificentService(url string, threshold int, ratio float64, timeout int, intervalTime int) *MagnificentService {
	return &MagnificentService{
		URL:          url,
		HTTPClient:   lib.GlobalHTTPClient,
		RespQueue:    []*CheckResult{},
		Threshold:    threshold,
		HealthRatio:  ratio,
		Timeout:      timeout,
		IntervalTime: intervalTime,
	}
}

func (s *MagnificentService) Check() *CheckResult {
	r := &CheckResult{
		ServiceName: "MagnificentService",
	}
	_, statusCode, err := s.HTTPClient.DispatchHTTPRequest(s.URL, lib.HTTPGet, nil, time.Duration(s.Timeout)*time.Second, 0)

	r.StatusCode = statusCode
	if err != nil {
		r.ErrorMsg = err.Error()
	}
	s.appendResult(r)
	return r
}

func (s *MagnificentService) appendResult(result *CheckResult) {
	s.RespQueue = append(s.RespQueue, result)
	if len(s.RespQueue) > s.Threshold {
		s.RespQueue = s.RespQueue[1:]
	}
}

func (s *MagnificentService) GetHealthStatus() string {
	success := 0
	for _, r := range s.RespQueue {
		if r.StatusCode == 200 {
			success++
		}
	}

	ratio := float64(success) / float64(len(s.RespQueue))
	if ratio == 1 {
		return "GREEN"
	}
	if ratio >= s.HealthRatio {
		return "YELLOW"
	}

	return "RED"
}
