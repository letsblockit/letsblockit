--- The two twitter templates have been deprecated
--- in the previous release, remove them from the DB

DELETE
FROM filter_instances
WHERE template_name IN ('twitter-tweets-by-hashtag', 'twitter-tweets-by-user');
