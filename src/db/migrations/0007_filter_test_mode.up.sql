ALTER TABLE filter_instances ADD COLUMN test_mode BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE user_preferences ADD COLUMN beta_features BOOLEAN NOT NULL DEFAULT FALSE;
