package health_check

import "services.order/internal/appservice/port"

func NewHealthCheckService(r port.Repo) *HealthCheckService {
	return &HealthCheckService{r}
}

type HealthCheckService struct {
	repo port.Repo
}

func (s *HealthCheckService) Check() error {
	return s.repo.Ping()
}
