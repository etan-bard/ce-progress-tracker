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

## Frontend UI

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
- **Styling**: CSS/TailwindCSS (to be determined)

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

### Component Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── ParticipantTable.vue      # Main table component
│   │   ├── CourseColumn.vue          # Course column component
│   │   └── SortControl.vue           # Sorting controls
│   ├── composables/
│   │   └── useParticipantData.js     # Data fetching logic
│   ├── App.vue                       # Main application
│   └── main.js                       # Entry point
├── public/                           # Static assets
└── package.json                      # Dependencies
```

### Data Flow

```
SQL Server → Backend API (GET /participant-courses) → Frontend Components → User Interface
```

### Sorting Implementation

The table supports sorting by:
- Participant ID (ascending/descending)
- Course completion percentage
- Last accessed date

### Error Handling

- Network errors display user-friendly messages
- Missing data shows "N/A" or appropriate placeholder
- Loading states during data fetch

### Future Enhancements

- Filtering by course or participant
- Export to CSV/Excel
- Visual progress indicators
- Detailed participant view

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