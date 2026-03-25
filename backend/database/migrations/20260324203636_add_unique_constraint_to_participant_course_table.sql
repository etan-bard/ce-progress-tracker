-- +goose Up
ALTER TABLE [ProgressTracker].[ParticipantToCourseMap]
    ADD CONSTRAINT UC_ParticipantToCourseMap_ParticipantId_CourseId UNIQUE (ParticipantId, CourseId);

-- +goose Down
ALTER TABLE [ProgressTracker].[ParticipantToCourseMap]
    DROP CONSTRAINT UC_ParticipantToCourseMap_ParticipantId_CourseId;
