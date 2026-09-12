-- Runs against the identity_db schema (created by the MySQL container /
-- provisioning, not by a migration).
CREATE TABLE IF NOT EXISTS users (
    id CHAR(36) NOT NULL,

    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,

    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,

    PRIMARY KEY (id),

    UNIQUE KEY ux_users_email (email),

    KEY ix_users_created_at (created_at)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
