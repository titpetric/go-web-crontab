CREATE TABLE IF NOT EXISTS `migrations` (
 `project` text,
 `filename` text,
 `statement_index` integer,
 `status` text,
 PRIMARY KEY (project, filename)
);

pragma journal_mode = WAL;
pragma synchronous = normal;
pragma temp_store = memory;
pragma mmap_size = 30000000000;
