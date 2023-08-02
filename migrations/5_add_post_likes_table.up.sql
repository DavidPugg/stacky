CREATE TABLE post_likes (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    post_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT FK_PostLikeUser FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT FK_PostLikePost FOREIGN KEY (post_id) REFERENCES posts(id),
    UNIQUE KEY unique_user_post (user_id, post_id)
);
