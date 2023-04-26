package service

import (
	"github.com/mohammadrabetian/ports/domain"
)

type PortService interface {
	AddOrUpdatePort(port *domain.Port) error
	GetPortByID(id string) (*domain.Port, error)
	ListPorts(page int, limit int) ([]*domain.Port, error)
}
