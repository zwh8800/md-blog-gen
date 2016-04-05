ALTER TABLE `mdblog`.`NoteTag`
ADD PRIMARY KEY (`note_id`, `tag_id`);

ALTER TABLE `mdblog`.`Note`
ADD INDEX `unique_id` (`unique_id` ASC),
ADD INDEX `removed` (`removed` ASC);

ALTER TABLE `mdblog`.`Tag`
CHANGE COLUMN `name` `name` VARCHAR(45) NOT NULL ,
ADD INDEX `name` (`name` ASC);
