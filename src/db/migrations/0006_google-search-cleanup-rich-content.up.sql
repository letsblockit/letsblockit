--- The only-results option has been deprecated, move instances to
--- the new 'rich-results' option as a partial replacement for it.

UPDATE filter_instances
SET params = filter_instances.params || jsonb_build_object('rich-results', true)
WHERE filter_name = 'google-search-cleanup' AND filter_instances.params -> 'only-results' = 'true';
