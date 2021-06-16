package health_check

func NewHealthCheckService() *HealthCheckService {
	return &HealthCheckService{}
}

type HealthCheckService struct {
}

func (s *HealthCheckService) Check() error {
	return nil
}
