ALTER TABLE follows DROP FOREIGN KEY FK_FollowFollower;
ALTER TABLE follows DROP FOREIGN KEY FK_FollowFollowee;

DROP TABLE follows;