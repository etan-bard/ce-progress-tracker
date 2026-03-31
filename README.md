# CE Progress Tracker

## CE Progress Sync Script
A script to export CE progress data from MongoDB, and patch the data into SQL Server Express.

## CE Progress Tracker Site
A web application to display CE progress data from SQL Server Express.

## API Implementation
The backend API is implemented in `/backend/main.go` and provides endpoints to interact with participant-course data stored in SQL Server Express. The API uses:
- Huma framework for OpenAPI documentation
- Chi router for HTTP routing
- SQLx for database operations
- MSSQL database service for data persistence

## Prerequisites
- Go v1.26
- Docker
- Optionally MongoDB and Sql Server Express

## Installation and Running
- For local development, a docker-file script is provided which will launch Mongo and SQL Server Express.
- Write all variables in `.env`, use `.env.example` as a template.
- Running the script by running: `go run main.go` In your CLI or IDE of choice.
- Optionally, you may add the flag `-courses=1,2,3` to specify which courses to migrate. This flag overwrites `.env`.

## Running the API
To run the API with Docker:
```bash
cd backend
docker-compose up api
```

This will start the API service along with all dependencies (MongoDB, SQL Server, and migrations). The API will be available on the port specified in your `.env` file (default: 8080).

## API Endpoints
- `GET /participant-courses` - Retrieves all participant-course mappings with completion and access date information

## OpenAPI Documentation
The API provides OpenAPI documentation at the following endpoint:
- `GET /docs` - OpenAPI documentation
- `GET /openapi.json` - OpenAPI specification in JSON format

## Running Migrations

When you run the API, migrations are automatically ran to update the database to the latest state. You may optionally migrate the database using the following commands.

**Up Migration(s):** 

`docker-compose run --rm goose up`

**Rollback (Once):** 

`docker-compose run --rm goose down`

**Create New Migration:**
`docker-compose run --rm goose create add_some_table sql`
