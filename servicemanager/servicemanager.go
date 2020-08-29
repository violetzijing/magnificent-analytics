package servicemanager

import (
	"fmt"
	"os"
	"sync"
	"time"

	"magnificent-analytics/services"

	log "github.com/sirupsen/logrus"
)

type ServiceManager struct {
	Services []services.Service
	F        *os.File
	Done     chan bool
	Mutex    sync.Mutex
	WG       sync.WaitGroup
}

func NewServiceManager(svs []services.Service) *ServiceManager {
	f, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return &ServiceManager{
		Services: svs,
		F:        f,
	}
}

func (m *ServiceManager) StartMonitor() {
	m.WG.Add(len(m.Services))
	for _, svc := range m.Services {
		go func(svc services.Service) {
			ticker := time.NewTicker(time.Duration(svc.GetIntervalTime()) * time.Second)
			for {
				select {
				case <-m.Done:
					m.WG.Done()
					return
				case <-ticker.C:
					result := svc.Check()
					healthyStatus := svc.GetHealthStatus()
					m.Log(result, healthyStatus)
				}
			}
		}(svc)
	}
}

func (m *ServiceManager) Log(result *services.CheckResult, status string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	str := fmt.Sprintf("Time: %v, Service: %s, StatusCode: %d", time.Now().UTC(), result.ServiceName, result.StatusCode)
	if result.ErrorMsg != "" {
		str = fmt.Sprintf("%s, ErrorMsg: %s", str, result.ErrorMsg)
	}
	str = fmt.Sprintf("%s, HealthStatus: %s \n", str, status)

	_, err := m.F.WriteString(str)
	if err != nil {
		log.Error("failed to log resuslt")
	}
}

func (m *ServiceManager) Stop() {
	m.F.Close()
	for i := 0; i < len(m.Services); i++ {
		m.Done <- true
	}
}
