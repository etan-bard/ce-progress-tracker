### Idempotent Script Application
This project utilizes an idempotent script for migrating data from MongoDB to SQL Server Express.

Run script with in `/backend/scripts`:
```bash
go run migrate_data.go
```

### Dockerized Development Environment
You are working within a Dockerized application. Reference  `/backend/docker-compose.yml` for the setup.
The following services are available:

- **mongodb**: Database service.
- **sqlserver**: Database service.
- **goose**: Database migration tool for SQL Server.
- **mongodb-restore**: Restores MongoDB from `mongo.dump` (manual profile).
- **mockery**: Mock generation tool (manual profile).
- **api**: Backend API service for participant-course data.

Run services with:
```bash
docker-compose up
```

For manual services (e.g., restore or mockery):
```bash
docker-compose --profile manual up mongodb-restore
```

## API Implementation
The backend API in `/backend/main.go` provides REST endpoints to interact with participant-course data. It uses:
- Huma framework for OpenAPI documentation
- SQLx for database interaction
- Chi router for HTTP routing
- MSSQL database service for data persistence

To run the API:
```bash
cd backend
docker-compose up api
```

The API will be available on the port specified in your `.env` file (default: 8080) and provides:
- `GET /participant-courses` - Retrieves all participant-course mappings with completion and access date information
