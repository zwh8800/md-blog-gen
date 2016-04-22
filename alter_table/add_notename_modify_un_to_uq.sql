ALTER TABLE `mdblog`.`Note`
  CHANGE COLUMN `id` `id` INT(10) NOT NULL AUTO_INCREMENT ,
  CHANGE COLUMN `unique_id` `unique_id` INT(10) NOT NULL ,
  ADD UNIQUE INDEX `unique_id_UNIQUE` (`unique_id` ASC);

ALTER TABLE `mdblog`.`Note`
  ADD COLUMN `notename` VARCHAR(80) NULL DEFAULT NULL AFTER `unique_id`,
  ADD UNIQUE INDEX `notename_UNIQUE` (`notename` ASC);

ALTER TABLE `mdblog`.`Note`
  ADD INDEX `notename` (`notename` ASC);
