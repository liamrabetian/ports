package repository

import (
	"errors"
	"fmt"

	"github.com/mohammadrabetian/ports/domain"
	"gorm.io/gorm"
)

type portMySQLRepository struct {
	db *gorm.DB
}

func NewPortMySQLRepository(db *gorm.DB) PortRepository {
	return &portMySQLRepository{db: db}
}

func (r *portMySQLRepository) createPort(port *domain.Port) error {
	if err := r.db.Create(port).Error; err != nil {
		return fmt.Errorf("failed to create port: %w", err)
	}
	return nil
}

func (r *portMySQLRepository) updatePort(port *domain.Port) error {
	if err := r.db.Model(port).Updates(port).Error; err != nil {
		return fmt.Errorf("failed to update port: %w", err)
	}
	return nil
}

/* NOTES: transactions maybe?! */
func (r *portMySQLRepository) AddOrUpdatePort(port *domain.Port) error {
	var existingPort domain.Port
	if err := r.db.Where("id = ?", port.ID).First(&existingPort).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.createPort(port)
		}
		return fmt.Errorf("failed to find port: %w", err)
	}
	return r.updatePort(port)
}

func (r *portMySQLRepository) FindByID(id string) (*domain.Port, error) {
	var port domain.Port
	if err := r.db.Where("id = ?", id).First(&port).Error; err != nil {
		return nil, err
	}
	return &port, nil
}

func (r *portMySQLRepository) ListPorts(page int, limit int) ([]*domain.Port, error) {
	var ports []*domain.Port

	offset := (page - 1) * limit
	result := r.db.Offset(offset).Limit(limit).Find(&ports)
	if result.Error != nil {
		return nil, result.Error
	}

	return ports, nil
}
