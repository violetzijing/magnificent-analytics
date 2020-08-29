package services

type Service interface {
	GetHealthStatus() string
	Check() *CheckResult
	GetIntervalTime() int
}

type CheckResult struct {
	ServiceName string
	StatusCode  int
	ErrorMsg    string
}
