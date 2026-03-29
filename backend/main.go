package main

import (
	"ce-progress-tracker/core"
	"ce-progress-tracker/database/mssql"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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

	// Add CORS middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

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
