--- Don't mark the search-results filter updated for the addition of the .onion parameters,
--- set them to false by default.

UPDATE filter_instances
SET params = filter_instances.params || jsonb_build_object('duckduckgo-onion', false, 'brave-onion', false)
WHERE template_name='search-results';
