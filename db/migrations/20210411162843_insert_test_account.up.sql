SET NAMES utf8 ;


LOCK TABLES `users` WRITE;
INSERT INTO `users` VALUES (1,CURDATE(),CURDATE(),NULL, 'Mr. Cool', '4', 'Boring','80 peace way', '20708', 'laurel','USA', '555-344-3434', '555-344-3434', '15');
UNLOCK TABLES;

LOCK TABLES `accounts` WRITE;
INSERT INTO `accounts` VALUES (1,CURDATE(),CURDATE(),NULL,'bin123@fake.com','$2a$14$xXeWMeclVoWO1SjJSz2kzOxUv8ZntbY6b7yzdNzPWzWGumskvAlF2', 1);
UNLOCK TABLES;

