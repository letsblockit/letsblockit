--- Folds instances of youtube-streams-chat into youtube-cleanup,
--- either by creating a new instance or updating the existing instance.
INSERT INTO filter_instances (user_id, filter_list_id, filter_name, params)
SELECT user_id, filter_list_id, 'youtube-cleanup', json_build_object('remove-stream-chat', true)
FROM filter_instances
WHERE filter_name = 'youtube-streams-chat'
ON CONFLICT (user_id, filter_name) DO UPDATE SET params    = filter_instances.params || excluded.params,
                                                 updated_at=now();

--- Remove obsolete instances
DELETE
FROM filter_instances
WHERE filter_name = 'youtube-streams-chat';
