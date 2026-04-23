CREATE TABLE IF NOT EXISTS users_backup (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    updated_at DATETIME,
    created_at DATETIME
);

INSERT INTO users_backup (id, username, email, password, updated_at, created_at)
SELECT id, username, email, password, updated_at, created_at
FROM users;

DROP TABLE users;
ALTER TABLE users_backup RENAME TO users;
