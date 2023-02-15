-- Allows instance operator to deviate from the default refresh period on lists with high traffic
ALTER TABLE filter_lists ADD COLUMN refresh_period_hours INT NULL;
