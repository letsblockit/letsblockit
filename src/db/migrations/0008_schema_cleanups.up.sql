-- Rename columns to better match template/filter semantics
ALTER TABLE filter_instances RENAME COLUMN filter_list_id TO list_id;
ALTER TABLE filter_instances RENAME COLUMN filter_name TO template_name;

-- Use downloaded_at IS NOT NULL instead
ALTER TABLE filter_lists DROP COLUMN downloaded;

-- Use timestampz everywhere
-- PG will convert existing values using the connection's TZ
ALTER TABLE banned_users ALTER created_at TYPE timestamptz;
ALTER TABLE banned_users ALTER lifted_at TYPE timestamptz;
ALTER TABLE filter_instances ALTER created_at TYPE timestamptz;
ALTER TABLE filter_instances ALTER updated_at TYPE timestamptz;
ALTER TABLE filter_lists ALTER created_at TYPE timestamptz;
ALTER TABLE filter_lists ALTER downloaded_at TYPE timestamptz;
ALTER TABLE user_preferences ALTER news_cursor TYPE timestamptz;
ALTER TABLE user_preferences ALTER news_cursor SET DEFAULT NOW();
