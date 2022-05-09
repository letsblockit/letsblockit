CREATE TABLE user_preferences
(
    user_id     uuid PRIMARY KEY,
    news_cursor timestamp NOT NULL DEFAULT (NOW() at time zone 'utc')
);

-- Populate news_cursor with the latest instance creation/modification timestamp for the user
INSERT INTO user_preferences (user_id, news_cursor)
SELECT user_id, GREATEST(MAX(created_at), MAX(updated_at))
FROM filter_instances
group by user_id;
