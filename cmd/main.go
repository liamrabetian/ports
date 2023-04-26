package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mohammadrabetian/ports/api"
	"github.com/mohammadrabetian/ports/domain"
	"github.com/mohammadrabetian/ports/handlers"
	"github.com/mohammadrabetian/ports/util"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		logrus.WithError(err).
			Fatal("cannot load config")
	}

	// Set up Logger
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// More readable logs for development env
	if config.Environment == "development" {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			DisableColors: false,
			FullTimestamp: true,
		})
		logrus.SetOutput(os.Stderr)
	}

	runGinServer(config)
}

func runGinServer(config util.Config) {
	server := api.NewServer(config)

	// Auto-migrate the database schema
	migrateDatabase(server.Store.SQL.DB)

	// Set up a channel to listen for signals
	signalChan := make(chan os.Signal, 1)
	// Notify the channel when an os.Interrupt or os.Kill signal is received
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		// Wait for a signal to be received
		sig := <-signalChan
		logrus.Warnf("Received signal: %v", sig)

		// Perform cleanup tasks, such as closing the database connection
		server.Store.SQL.Close()
		logrus.Warn("Database connection closed")

		os.Exit(0)
	}()

	/* NOTES: cli */
	err := handlers.ProcessJSONFile(config.FilePath)
	if err != nil {
		logrus.WithError(err).
			Fatal("error processing json file")
	} else {
		logrus.Info("JSON file processed successfuly")
	}

	err = server.Start(config.HTTPServer.Address)
	if err != nil {
		logrus.WithError(err).
			Fatal("cannot run server")
	}
}

func migrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(&domain.Port{})
	if err != nil {
		logrus.Fatalf("Failed to migrate database: %v", err)
	}
	logrus.Info("Database migration completed")
}
