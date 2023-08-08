ALTER TABLE posts
ADD COLUMN user_id INT NOT NULL DEFAULT 0 AFTER id;

ALTER TABLE posts
ADD KEY user_id_idx (user_id);