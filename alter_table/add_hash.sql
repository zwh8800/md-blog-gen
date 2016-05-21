ALTER TABLE `mdblog`.`Note`
  ADD COLUMN `hash` VARCHAR(32) NOT NULL AFTER `timestamp`;
