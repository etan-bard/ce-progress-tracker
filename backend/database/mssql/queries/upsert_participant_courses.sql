MERGE INTO %s AS Target
USING (
    SELECT 
        CAST(v.ParticipantId AS INT), 
        CAST(v.CourseId AS INT), 
        CAST(v.DateLastAccessed AS DATETIME2),
        CAST(v.CourseCompletion AS REAL)
    FROM (VALUES %s) AS v (ParticipantId, CourseId, DateLastAccessed, CourseCompletion)
) AS Source (ParticipantId, CourseId, DateLastAccessed, CourseCompletion)
ON Target.ParticipantId = Source.ParticipantId AND Target.CourseId = Source.CourseId
WHEN MATCHED AND (
    Target.DateLastAccessed IS DISTINCT FROM Source.DateLastAccessed 
    OR Target.CourseCompletion IS DISTINCT FROM Source.CourseCompletion) THEN
    UPDATE SET 
        Target.DateLastAccessed = Source.DateLastAccessed,
        Target.CourseCompletion = Source.CourseCompletion
WHEN NOT MATCHED BY TARGET THEN
    INSERT (ParticipantId, CourseId, DateLastAccessed, CourseCompletion)
    VALUES (Source.ParticipantId, Source.CourseId, Source.DateLastAccessed, Source.CourseCompletion)
OUTPUT $action AS Action;