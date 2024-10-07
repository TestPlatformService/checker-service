CREATE TABLE IF NOT EXISTS submited (
    id UUID PRIMARY KEY,
    code TEXT,
    user_task_id UUID ,
    submited_at TIMESTAMP
);
