package main

import (
	"magnificent-analytics/lib"
	"magnificent-analytics/servicemanager"
	"magnificent-analytics/services"
)

func main() {
	registerMonitor()
}

func registerMonitor() {
	cfg := lib.ParseConfig()
	magnificentSVC := services.NewMagnificentService(cfg.URL, cfg.Threshold, cfg.HealthRatio, cfg.Timeout, cfg.IntervalTime)
	manager := servicemanager.NewServiceManager([]services.Service{magnificentSVC})
	manager.StartMonitor()
	manager.WG.Wait()
}
