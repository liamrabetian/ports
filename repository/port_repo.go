package repository

import (
	"github.com/mohammadrabetian/ports/domain"
)

type PortRepository interface {
	AddOrUpdatePort(port *domain.Port) error
	FindByID(id string) (*domain.Port, error)
	ListPorts(page int, limit int) ([]*domain.Port, error)
	updatePort(port *domain.Port) error
	createPort(port *domain.Port) error
}
