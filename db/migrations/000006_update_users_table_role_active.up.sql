CREATE TABLE IF NOT EXISTS users_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user')),
    active BOOLEAN NOT NULL DEFAULT 1,
    updated_at DATETIME,
    created_at DATETIME
);

INSERT INTO users_new (id, username, email, password, role, active, updated_at, created_at)
SELECT id, username, email, password, 'user', 1, updated_at, created_at
FROM users;

DROP TABLE users;
ALTER TABLE users_new RENAME TO users;
