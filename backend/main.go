package main

import (
	"ce-progress-tracker/core"
	"ce-progress-tracker/database/mssql"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	// Initialize Viper to load .env
	config, err := core.NewConfig()
	if err != nil {
		log.Fatal("Failed to load config", zap.Error(err))
	}

	// Initialize Logger with configured level
	logger, err := core.MakeLogger(config.GetLogLevel())
	if err != nil {
		log.Fatal("Failed to create logger", zap.Error(err))
	}
	logger.Info("Configuration and logger initialized", zap.String("level", string(config.GetLogLevel())))

	// Initialize database
	db, err := mssql.NewMSSQLDatabaseService(
		config.GetMSSQLConnectionString(),
	)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer db.Close()
	logger.Info("Database initialized")

	router := chi.NewRouter()
	humaConfig := huma.DefaultConfig("CE Progress Tracker", "1.0.0")
	api := humachi.New(router, humaConfig)

	participantCourseRepository := mssql.NewParticipantCourseRepository(db)
	participantCourseController := mssql.NewParticipantCourseController(participantCourseRepository)
	participantCourseController.Register(api)

	logger.Debug("API routes registered")
	port := config.GetAPIPort()
	logger.Debug("Starting server on port", zap.String("port", port))
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
