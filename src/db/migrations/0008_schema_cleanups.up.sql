-- Rename columns to better match template/filter semantics
alter table filter_instances rename column filter_list_id to list_id;
alter table filter_instances rename column filter_name to template_name;

-- Use downloaded_at IS NOT NULL instead
alter table filter_lists drop column downloaded;