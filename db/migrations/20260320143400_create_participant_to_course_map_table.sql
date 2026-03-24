-- +goose Up
-- Create schema if it does not exist
IF SCHEMA_ID('progress_tracker') IS NULL
    EXEC('CREATE SCHEMA progress_tracker');

-- Create table if it does not exist
IF OBJECT_ID('[progress_tracker].[participant_to_course_map]', 'U') IS NULL
CREATE TABLE [progress_tracker].[participant_to_course_map] (
    [Id] INT IDENTITY(1,1) PRIMARY KEY,
    [ParticipantId] INT NOT NULL,
    [CourseId] INT NOT NULL,
    [LastAccessedAt] DATETIME2 NULL,
    [CourseCompletion] BIT NOT NULL DEFAULT(0)
    );

-- +goose Down
-- Drop table if it exists
IF OBJECT_ID('[progress_tracker].[participant_to_course_map]', 'U') IS NOT NULL
DROP TABLE [progress_tracker].[participant_to_course_map];