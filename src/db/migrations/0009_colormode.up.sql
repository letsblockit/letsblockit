CREATE TYPE color_mode AS ENUM ('auto', 'dark', 'light');
ALTER TABLE user_preferences ADD COLUMN color_mode color_mode NOT NULL DEFAULT 'auto';
