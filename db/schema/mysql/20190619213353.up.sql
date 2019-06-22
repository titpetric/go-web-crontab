ALTER TABLE `logs` ADD `output` JSON NOT NULL AFTER `name`, ADD `exit_code` INT NOT NULL AFTER `output`;
