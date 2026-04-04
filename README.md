# CE Progress Tracker
This mono-repository contains 3 components for the assignment. Task 1 includes a script to migrate data from MongoDB to SQL Server Express. 
Task 2 includes a backend API and a frontend website to display the migrated data.

## Task 1: CE Progress Sync Script
A script to export CE progress data from MongoDB, and patch the data into SQL Server Express.

## CE Progress Tracker Backend API
The backend API is implemented in `/backend/main.go` and provides endpoints to interact with participant-course data stored in SQL Server Express. The API uses:
- Huma framework for OpenAPI documentation
- Chi router for HTTP routing
- SQLx for database operations against SQL Server Express

## Task 2: CE Progress Tracker Site
The frontend application is located in the `/frontend` directory and provides a user interface to view participant progress data from SQL Server.

### UI Requirements

- **Cross-tab format display**
  - Participant IDs displayed as rows
  - Courses displayed as columns
  - Each course column shows:
    - Completion percentage
    - Last accessed date

- **Sortable table**
  - Users can sort by participant ID
  - Users can sort by course-related columns

- **Data handling**
  - Reads data from backend API endpoint
  - Handles missing data gracefully

### Technology Stack

- **Framework**: Vue.js
- **Components**: Vue Single File Components
- **State Management**: Vue Composition API
- **Styling**: Vuetify Framework

### Acceptance Criteria

1. ✅ UI reads data from backend API endpoint
2. ✅ Participant IDs displayed vertically as rows
3. ✅ Courses displayed horizontally as columns
4. ✅ Each participant/course intersection shows:
   - Completion percentage
   - Last accessed date
5. ✅ Table is sortable by participant ID
6. ✅ Table is sortable by course-related columns
7. ✅ Missing data handled gracefully with appropriate UI feedback

### Development Setup

```bash
cd frontend
npm install
npm run dev
```

### Building for Production

```bash
cd frontend
npm run build
```

### API Integration

The frontend connects to the backend API at `GET /participant-courses` endpoint to fetch participant progress data. This endpoint is implemented in `backend/database/mssql/participant_course_controller.go` and returns all participant-course mappings with completion and access date information.

### Sorting Implementation

The table supports sorting by:
- Participant ID (ascending/descending)
- Course completion percentage
- Last accessed date

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
