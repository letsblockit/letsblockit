--- The rules for channels were previously enabled by the subscriptions option,
--- enable it if needed for continuity.

UPDATE filter_instances
SET params = filter_instances.params || jsonb_build_object('channels', true)
WHERE template_name='youtube-watched' AND filter_instances.params -> 'subscriptions' = 'true';
