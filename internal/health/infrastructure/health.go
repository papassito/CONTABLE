package infrastructure

type HealthChecker interface {
	Ping() bool
}

type healthChecker struct{}

func NewHealthChecker() HealthChecker {
	return &healthChecker{}
}

func (h *healthChecker) Ping() bool {
	return true
}
