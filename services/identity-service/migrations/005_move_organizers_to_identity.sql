USE identity_db;

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
    UNIQUE KEY ux_organizers_slug (slug),
    KEY ix_organizers_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO identity_db.organizers (id, email, name, phone, slug, is_active, created_at, updated_at)
SELECT id, email, name, phone, slug, is_active, created_at, updated_at
FROM ticket_db.organizers
ON DUPLICATE KEY UPDATE
    email = VALUES(email),
    name = VALUES(name),
    phone = VALUES(phone),
    slug = VALUES(slug),
    is_active = VALUES(is_active),
    updated_at = VALUES(updated_at);

SET @has_membership_fk = (
    SELECT COUNT(*)
    FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = 'identity_db'
      AND TABLE_NAME = 'users'
      AND CONSTRAINT_NAME = 'fk_users_organizer'
      AND CONSTRAINT_TYPE = 'FOREIGN KEY'
);
SET @add_membership_fk = IF(
    @has_membership_fk = 0,
    'ALTER TABLE identity_db.users ADD CONSTRAINT fk_users_organizer FOREIGN KEY (organizer_id) REFERENCES identity_db.organizers(id) ON DELETE SET NULL',
    'SELECT 1'
);
PREPARE statement FROM @add_membership_fk;
EXECUTE statement;
DEALLOCATE PREPARE statement;
