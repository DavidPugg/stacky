CREATE TABLE follows (
    id INT NOT NULL AUTO_INCREMENT,
    follower_id INT NOT NULL,
    followee_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY follower_id_idx (follower_id),
    KEY followee_id_idx (followee_id),
    UNIQUE KEY unique_follower_followee (follower_id, followee_id)
);