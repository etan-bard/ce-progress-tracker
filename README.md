# CE Progress Tracker

## CE Progress Sync Script
A script to export CE progress data from MongoDB, and patch the data into SQL Server Express.

## CE Progress Tracker Site
A web application to display CE progress data from SQL Server Express.

## API Implementation
The backend API is implemented in `/backend/main.go` and provides endpoints to interact with participant-course data stored in SQL Server Express. The API uses:
- Huma framework for OpenAPI documentation
- Chi router for HTTP routing
- SQLC for database operations
- MSSQL database service for data persistence

## Prerequisites
- Go v1.26
- MongoDB
- SQL Server Express

## Installation and Running
- Ensure you have the DB credentials and connection strings set up for MongoDB and SQL Server Express.
- Running the script by running: `go run main.go` In your CLI.

## Running the API
To run the API with Docker:
```bash
cd backend
docker-compose up api
```

This will start the API service along with all dependencies (MongoDB, SQL Server, and migrations). The API will be available on the port specified in your `.env` file (default: 8080).

## API Endpoints
- `GET /participant-courses` - Retrieves all participant-course mappings with completion and access date information

## Plan

### Containerization (Development and Environment Setup)
- Dockerized MongoDB and SQL Server Express for easy development ease and demoing in a sandbox environment
- Dockerized Goose for SQL Server Express migrations
- Dockerized Goose for generation of schema.sql to be ingested by SQLc

### Data Interface Layer
- SQLc to Interface with SQL Server Express
- Use Mongo Go driver for MongoDB

### Script Outline
- Pull the Course List from `.env`
- Seed the SQL Server Express table from MongoDB
  - Prevent duplicate records
  - Insert new records (if any)
  - Update existing records (if any)
- When script is finished, log the number of rows that were inserted and updated

### Restrictions/Constraints
- Script must be idempotent: subsequent runs with identical MongoDB and course list states should not modify the SQL Server Express DB
- Script logging should be less verbose

## Running Migrations

**Up Migration(s):** 

`docker-compose run --rm goose up`

**Rollback (Once):** 

`docker-compose run --rm goose down`

**Create New Migration:**
`docker-compose run --rm goose create add_some_table sql`