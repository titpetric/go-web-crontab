CREATE TABLE `jobs` (
    `name` TEXT NOT NULL PRIMARY KEY,
    `description` TEXT NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    `deleted_at` TIMESTAMP NOT NULL
);

CREATE TABLE `logs` (
    `name` TEXT NOT NULL,
    `stamp` TIMESTAMP NOT NULL,
    `duration` INTEGER NOT NULL,
    `output` TEXT NOT NULL DEFAULT '',
    `exit_code` INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (`name`, `stamp`)
);
