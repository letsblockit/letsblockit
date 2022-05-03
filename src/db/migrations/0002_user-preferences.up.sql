CREATE TABLE user_preferences
(
    user_id     uuid PRIMARY KEY,
    latest_news timestamp NOT NULL DEFAULT NOW()
);

-- Populate latest_news with the latest instance creation/modification timestamp for the user
INSERT INTO user_preferences (user_id, latest_news)
SELECT user_id, GREATEST(MAX(created_at), MAX(updated_at))
FROM filter_instances
group by user_id;
