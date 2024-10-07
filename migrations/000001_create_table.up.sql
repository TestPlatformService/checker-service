CREATE TABLE IF NOT EXISTS submitted (
    id UUID PRIMARY KEY,
    code TEXT,
    user_task_id UUID,
    submitted_at TIMESTAMP
);
