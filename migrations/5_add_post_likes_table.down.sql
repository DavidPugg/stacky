ALTER TABLE post_likes DROP FOREIGN KEY FK_PostLikeUser;
ALTER TABLE post_likes DROP FOREIGN KEY FK_PostLikePost;

DROP TABLE post_likes;
