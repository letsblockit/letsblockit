CREATE TABLE filter_lists
(
    id         SERIAL PRIMARY KEY,
    user_id    uuid      NOT NULL,
    token      uuid      NOT NULL DEFAULT gen_random_uuid(),
    created_at timestamp NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_lists_by_token ON filter_lists USING btree (token);
CREATE UNIQUE INDEX idx_lists_by_user ON filter_lists USING btree (user_id);
ALTER TABLE filter_lists
    ADD COLUMN downloaded bool NOT NULL DEFAULT FALSE;
ALTER TABLE filter_lists
    ADD COLUMN downloaded_at timestamp;

CREATE TABLE filter_instances
(
    id             SERIAL PRIMARY KEY,
    user_id        uuid      NOT NULL,
    filter_list_id INTEGER   NOT NULL REFERENCES filter_lists (id) ON DELETE CASCADE,
    filter_name    text      NOT NULL,
    params         jsonb,
    created_at     timestamp NOT NULL DEFAULT NOW(),
    updated_at     timestamp
);

CREATE INDEX idx_instances_by_list ON filter_instances USING btree (filter_list_id);
CREATE UNIQUE INDEX idx_instances_by_user_and_filter ON filter_instances USING btree (user_id, filter_name);

CREATE TABLE banned_users
(
    id          SERIAL PRIMARY KEY,
    user_id     uuid      NOT NULL,
    created_at  timestamp NOT NULL DEFAULT NOW(),
    reason      text      NOT NULL,
    lifted_at   timestamp,
    lift_reason text
);
