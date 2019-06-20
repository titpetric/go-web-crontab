CREATE TABLE jobs (
    name        varchar(64) NOT NULL,
    description varchar(255) NOT NULL,
    PRIMARY KEY (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE logs (
    name     varchar(64) NOT NULL,
    success  bool NOT NULL,
    output   longtext NOT NULL,
    stamp    datetime NOT NULL,
    duration bigint(20) NOT NULL,
    PRIMARY KEY (name,stamp)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

