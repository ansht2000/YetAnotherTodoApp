-- +goose Up
CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    body TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT false,
    goal_id UUID NOT NULL,
    CONSTRAINT fk_goal_id
        FOREIGN KEY (goal_id)
        REFERENCES goals(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE tasks;