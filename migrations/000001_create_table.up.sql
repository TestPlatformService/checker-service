CREATE TABLE IF NOT EXISTS submited (
    id UUID PRIMARY KEY,
    code TEXT,
    user_task_id UUID REFERENCES user_tasks(id) ON DELETE CASCADE,
    submited_at TIMESTAMP
);
