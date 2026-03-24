package core

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all the configuration for the application
type Config struct {
	Viper   *viper.Viper
	FlagSet *flag.FlagSet
	courses *string
}

// NewConfig creates a new Config instance that is DI ready
func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigType("env")
	v.SetConfigName(".env")
	v.AddConfigPath(".")
	v.AutomaticEnv() // Fallback to system environment variables

	if err := v.ReadInConfig(); err != nil {
		// Only return error if .env file exists but failed to read
		if _, statErr := os.Stat(".env"); statErr == nil {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	ret := &Config{
		Viper:   v,
		FlagSet: flag.CommandLine,
	}

	// Establishes the command line flags, then parses them
	ret.FlagSet.StringVar(ret.courses, "courses", "", "comma-separated List of course IDs to process")
	err := ret.FlagSet.Parse(os.Args[1:])
	if err != nil {
		return nil, fmt.Errorf("error parsing command line flags: %w", err)
	}
	return ret, nil
}

func (c *Config) GetCourseIDs() *[]int {
	if c.Viper.IsSet("COURSE_IDS") {
		courseIDs := c.Viper.GetIntSlice("COURSE_IDS")
		return &courseIDs
	} else if c.courses != nil {
		courseStrings := strings.Split(*c.courses, ",")
		courseIDs := make([]int, 0, len(*c.courses))
		for _, s := range courseStrings {
			id, err := strconv.Atoi(strings.TrimSpace(s))
			if err != nil {
				log.Fatalf("Invalid course ID: %s", s)
			}
			courseIDs = append(courseIDs, id)
		}
		return &courseIDs
	}
	return nil
}

// GetLogLevel returns the LOG_LEVEL from the configuration
func (c *Config) GetLogLevel() LogLevel {
	return LogLevel(c.getRequiredValue("LOG_LEVEL"))
}

// GetMSSQLConnectionString returns the MSSQL connection string from the configuration
func (c *Config) GetMSSQLConnectionString() string {
	user := c.getRequiredValue("MSSQL_USER")
	password := c.getRequiredValue("MSSQL_SA_PASSWORD")
	host := c.getRequiredValue("MSSQL_HOST")
	port := c.getRequiredValue("MSSQL_PORT")
	return fmt.Sprintf("sqlserver://%s:%s@%s:%s?encrypt=disable", user, password, host, port)
}

// GetMongoURI returns the MongoDB connection URI from the configuration
func (c *Config) GetMongoURI() string {
	host := c.getRequiredValue("MONGO_HOST")
	port := c.getRequiredValue("MONGO_PORT")
	return fmt.Sprintf("mongodb://%s:%s", host, port)
}

// GetMongoDBName returns the MongoDB database name from the configuration
func (c *Config) GetMongoDBName() string {
	return c.getRequiredValue("MONGO_DB_NAME")
}

func (c *Config) GetScriptBatchSize() int {
	batchSize := c.Viper.GetInt("SCRIPT_BATCH_SIZE")
	if batchSize <= 0 || batchSize > 1000 {
		batchSize = 1000
	}
	return batchSize
}

// getRequiredValue returns value with priority: .env key > env var.
// If not found, logs an error.
func (c *Config) getRequiredValue(envKey string) string {
	if c.Viper.IsSet(envKey) {
		return c.Viper.GetString(envKey)
	}
	log.Fatalf("Missing required configuration: %s", envKey)
	return ""
}
