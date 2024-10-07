CREATE TABLE IF NOT EXISTS submitted (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT,
    user_task_id UUID,
    submitted_at TIMESTAMP
);
