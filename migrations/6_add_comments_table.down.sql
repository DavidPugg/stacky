ALTER TABLE comments DROP KEY post_id_idx;
ALTER TABLE comments DROP KEY user_id_idx;

DROP TABLE comments;