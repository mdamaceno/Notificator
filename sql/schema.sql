CREATE TABLE IF NOT EXISTS messages (
    id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    service varchar(255) NOT NULL,
    title varchar(255) NOT NULL,
    body text NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS destinations (
    id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    message_id uuid NOT NULL,
    receiver varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE
);
