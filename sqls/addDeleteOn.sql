ALTER TABLE `blog_tag` ADD `deleted_on` int(10) unsigned;
ALTER TABLE `blog_article` ADD `deleted_on` int(10) unsigned;


ALTER TABLE `blog_tag` MODIFY `deleted_on` int(10) unsigned not NULL DEFAULT 0;
ALTER TABLE `blog_article` MODIFY `deleted_on` int(10) unsigned not NULL DEFAULT 0;