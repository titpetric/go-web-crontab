ALTER TABLE logs ADD (
	success  bool NOT NULL,
    output   longtext NOT NULL
);
