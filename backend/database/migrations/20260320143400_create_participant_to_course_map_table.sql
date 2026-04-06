-- +goose Up
-- Create schema if it does not exist
IF SCHEMA_ID('ProgressTracker') IS NULL
    EXEC('CREATE SCHEMA ProgressTracker');

-- Create table if it does not exist
IF OBJECT_ID('[ProgressTracker].[ParticipantToCourseMap]', 'U') IS NULL
CREATE TABLE [ProgressTracker].[ParticipantToCourseMap] (
    [Id] INT IDENTITY(1,1) PRIMARY KEY,
    [ParticipantId] INT NOT NULL,
    [CourseId] INT NOT NULL,
    [DateFirstAccessed] DATETIME2 NULL,
    [DateLastAccessed] DATETIME2 NULL,
    [CourseCompletion] REAL NOT NULL DEFAULT(0)
    );

-- +goose Down
-- Drop table if it exists
IF OBJECT_ID('[ProgressTracker].[ParticipantToCourseMap]', 'U') IS NOT NULL
DROP TABLE [ProgressTracker].[ParticipantToCourseMap];