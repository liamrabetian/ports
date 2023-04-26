package repository_test

import (
	"testing"

	"github.com/mohammadrabetian/ports/domain"
	"github.com/mohammadrabetian/ports/pkg/mysql"
	"github.com/mohammadrabetian/ports/repository"
	"github.com/mohammadrabetian/ports/util"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var testPort = &domain.Port{
	ID:      "AEAJM",
	Name:    "Jebel Ali",
	City:    "Jebel Ali",
	Country: "United Arab Emirates",
}

func setupDatabase() *gorm.DB {
	config, err := util.LoadTestConfig("../.", "ports_test")
	if err != nil {
		panic("failed to load test config")
	}

	db := mysql.NewDatabase(config)

	err = db.DB.AutoMigrate(&domain.Port{})
	if err != nil {
		panic("failed to migrate the test database schema")
	}

	return db.DB
}

func TestPortMySQLRepositoryIntegration(t *testing.T) {
	db := setupDatabase()
	repo := repository.NewPortMySQLRepository(db)

	// Cleanup: delete the test port after each test
	t.Cleanup(func() {
		db.Delete(&domain.Port{}, testPort.ID)
	})
	t.Run("AddOrUpdatePort", func(t *testing.T) {
		err := repo.AddOrUpdatePort(testPort)
		assert.Nil(t, err)
	})

	t.Run("FindByID", func(t *testing.T) {
		foundPort, err := repo.FindByID(testPort.ID)
		assert.Nil(t, err)
		assert.Equal(t, testPort.ID, foundPort.ID)
		assert.Equal(t, testPort.Name, foundPort.Name)
		assert.Equal(t, testPort.City, foundPort.City)
		assert.Equal(t, testPort.Country, foundPort.Country)
	})

	t.Run("ListPorts", func(t *testing.T) {
		ports, err := repo.ListPorts(1, 10)
		assert.Nil(t, err)
		assert.GreaterOrEqual(t, len(ports), 1)
	})
}
