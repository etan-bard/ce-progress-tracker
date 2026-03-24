package main

import (
	"ce-progress-tracker/internal/core"
	"ce-progress-tracker/internal/database/mongo"
	"ce-progress-tracker/internal/database/mssql"
	"ce-progress-tracker/internal/services"
	"context"
	"log"

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
	dsn := config.GetMSSQLConnectionString()
	db, err := mssql.NewMSSQLDatabaseService(dsn)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer db.Close()
	logger.Info("Database initialized")

	// Initialize MongoDB service
	mongoService, err := mongo.NewMongoDBService(context.Background(), config.GetMongoURI(), config.GetMongoDBName())
	if err != nil {
		logger.Fatal("Failed to initialize MongoDB service", zap.Error(err))
	}
	defer mongoService.Close(context.Background())
	logger.Info("MongoDB service initialized")

	participantCourseRepository := mssql.NewParticipantCourseRepository(db)
	takesAnonymizedRepository := mongo.NewTakesAnonymizedRepository(mongoService)
	mapper := services.NewParticipantCourseMapper()
	batchSize := config.GetScriptBatchSize()

	dataMigrationStrategy := services.NewBatchDataMigrationStrategy(mapper, participantCourseRepository, takesAnonymizedRepository, logger, batchSize)
	inserted, updated, err := dataMigrationStrategy.Execute(context.Background(), nil)

	logger.Info("Migration completed", zap.Int("total_inserted", inserted), zap.Int("total_updated", updated))
}
