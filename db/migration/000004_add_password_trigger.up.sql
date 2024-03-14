-- +migrate Up
CREATE
OR REPLACE FUNCTION update_password_changed_at() RETURNS TRIGGER AS $$ BEGIN 

  -- Check if the hashed_password has changed
  IF OLD.hashed_password IS DISTINCT
  FROM
    NEW.hashed_password THEN -- Update password_changed_at only if the password has changed
    NEW.password_changed_at := NOW();

  END IF;

RETURN NEW;
END;

$$ LANGUAGE plpgsql;

CREATE TRIGGER update_password_trigger BEFORE
UPDATE
  OF hashed_password ON users FOR EACH ROW EXECUTE FUNCTION update_password_changed_at();