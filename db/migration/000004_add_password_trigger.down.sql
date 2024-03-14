-- +migrate Down
DROP TRIGGER update_password_trigger ON users;

DROP FUNCTION IF EXISTS update_password_changed_at();
