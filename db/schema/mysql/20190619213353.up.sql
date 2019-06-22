ALTER TABLE logs ADD (
	exit_code  int NOT NULL,
    output     json NOT NULL
);
