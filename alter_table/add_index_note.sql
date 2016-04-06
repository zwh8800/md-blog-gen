ALTER TABLE `mdblog`.`NoteTag`
ADD PRIMARY KEY (`note_id`, `tag_id`);

ALTER TABLE `mdblog`.`Note`
ADD INDEX `unique_id` (`unique_id` ASC),
ADD INDEX `removed` (`removed` ASC);

ALTER TABLE `mdblog`.`Tag`
CHANGE COLUMN `name` `name` VARCHAR(45) NOT NULL ,
ADD INDEX `name` (`name` ASC);

ALTER TABLE `mdblog`.`Note`
CHANGE COLUMN `title` `title` varchar(240) NOT NULL,
CHANGE COLUMN `url` `url` varchar(240) NOT NULL,
CHANGE COLUMN `content` `content` mediumtext NOT NULL,
CHANGE COLUMN `timestamp` `timestamp` datetime NOT NULL,
CHANGE COLUMN `removed` `removed` tinyint(1) unsigned NOT NULL DEFAULT '0';
