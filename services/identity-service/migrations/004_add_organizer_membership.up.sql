ALTER TABLE users ADD COLUMN organizer_id CHAR(36) NULL AFTER role;
ALTER TABLE users ADD KEY ix_users_organizer_id (organizer_id);
