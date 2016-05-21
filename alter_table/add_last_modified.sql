ALTER TABLE `mdblog`.`Note` 
  ADD COLUMN `last_modified` DATETIME NOT NULL AFTER `timestamp`;
