package service

import (
	"github.com/mohammadrabetian/ports/domain"
	"github.com/mohammadrabetian/ports/repository"
)

type portService struct {
	repo repository.PortRepository
}

func NewPortService(repo repository.PortRepository) PortService {
	return &portService{repo: repo}
}

func (s *portService) GetPortByID(id string) (*domain.Port, error) {
	return s.repo.FindByID(id)
}

func (s *portService) AddOrUpdatePort(port *domain.Port) error {
	return s.repo.AddOrUpdatePort(port)
}

func (s *portService) ListPorts(page int, limit int) ([]*domain.Port, error) {
	return s.repo.ListPorts(page, limit)
}
