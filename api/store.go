package api

import (
	"github.com/mohammadrabetian/ports/handlers"
	"github.com/mohammadrabetian/ports/pkg/mysql"
	"github.com/mohammadrabetian/ports/repository"
	"github.com/mohammadrabetian/ports/service"
	"github.com/mohammadrabetian/ports/util"
)

// All stores e.g. nosql,sql
type Store struct {
	SQL mysql.SQLDatabase
}

func NewStore(config util.Config) *Store {
	db := mysql.NewDatabase(config)

	// initialize the repo
	portRepo := repository.NewPortMySQLRepository(db.DB)

	// initialize the service and handlers
	portSvc := service.NewPortService(portRepo)
	handlers.InitPortHandlers(portSvc)

	return &Store{
		SQL: db,
	}
}
