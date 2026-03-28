### Idempotent Script Application
This project utilizes an idempotent script for migrating data from MongoDB to SQL Server Express.

Run script with:
```bash
go run main.go
```

### Dockerized Development Environment
You are working within a Dockerized application utilizing `docker-compose.yml`. The following services are available:

- **mongodb**: Database service on port `27017`.
- **sqlserver**: Database service on port `1433`.
- **goose**: Database migration tool for SQL Server.
- **mongodb-restore**: Restores MongoDB from `mongo.dump` (manual profile).
- **mockery**: Mock generation tool (manual profile).

Run services with:
```bash
docker-compose up
```

For manual services (e.g., restore or mockery):
```bash
docker-compose --profile manual up mongodb-restore
```