ALTER TABLE post_likes DROP KEY user_id_idx;
ALTER TABLE post_likes DROP KEY post_id_idx;

DROP TABLE post_likes;
