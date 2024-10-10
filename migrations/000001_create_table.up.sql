CREATE TYPE status_enum AS ENUM ('Accepted', 'Wrong answer', 'Compilation error');

CREATE TABLE IF NOT EXISTS submitted (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_name VARCHAR(200),
    status status_enum,
    lang VARCHAR(50),
    compiled_time INT,
    compiled_memory INT,
    code TEXT,
    user_task_id UUID,
    submitted_at TIMESTAMP
);
