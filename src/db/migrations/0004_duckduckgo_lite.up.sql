--- Don't mark the search-results filter updated for users who disabled duckduckgo support (see PR 167)
--- This is achieved by updating the params to add the corresponding fields, set to false.
--- Users who enabled DDG support will see the filter as updated, and offered the new options.

UPDATE filter_instances
SET params = filter_instances.params || jsonb_build_object('duckduckgo-lite', false, 'duckduckgo-html', false)
WHERE filter_name='search-results' AND filter_instances.params -> 'duckduckgo' = 'false';
