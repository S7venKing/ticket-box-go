CREATE DATABASE IF NOT EXISTS ticket_db
  CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE ticket_db;

CREATE TABLE IF NOT EXISTS organizers (
    id CHAR(36) NOT NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NULL,
    slug VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY ux_organizers_email (email),
    UNIQUE KEY ux_organizers_slug (slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS events (
    id CHAR(36) NOT NULL,
    organizer_id CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NULL,
    venue VARCHAR(255) NOT NULL,
    start_at DATETIME(6) NOT NULL,
    end_at DATETIME(6) NOT NULL,
    capacity INT NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'draft',
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    KEY ix_events_organizer_created (organizer_id, created_at),
    KEY ix_events_status_created (status, created_at),
    CONSTRAINT fk_events_organizer FOREIGN KEY (organizer_id) REFERENCES organizers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
