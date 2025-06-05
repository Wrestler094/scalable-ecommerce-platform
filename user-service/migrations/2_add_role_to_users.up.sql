ALTER TABLE users
ADD COLUMN role TEXT NOT NULL DEFAULT 'user',
ADD CONSTRAINT check_user_role CHECK (role IN ('user', 'admin'));
