-- +goose Up
CREATE TABLE goals (
    id UUID PRIMARY KEY,
    body TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT false,
    user_id UUID NOT NULL,
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE goals;