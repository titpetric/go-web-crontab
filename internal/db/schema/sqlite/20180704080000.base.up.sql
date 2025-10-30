CREATE TABLE `jobs` (
    `name` TEXT NOT NULL PRIMARY KEY,
    `description` TEXT NOT NULL
);

CREATE TABLE `logs` (
    `name` TEXT NOT NULL,
    `stamp` TEXT NOT NULL,
    `duration` INTEGER NOT NULL,
    PRIMARY KEY (`name`, `stamp`)
);
