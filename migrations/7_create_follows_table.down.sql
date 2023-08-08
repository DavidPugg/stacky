ALTER TABLE follows DROP KEY follower_id_idx;
ALTER TABLE follows DROP KEY followee_id_idx;

DROP TABLE follows;