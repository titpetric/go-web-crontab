// Code generated by statik. DO NOT EDIT.

// Package contains static assets.
package mysql

var	Asset = "PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x1a\x00	\x0020180704080000.base.up.sqlUT\x05\x00\x01\x80Cm8CREATE TABLE `jobs` (\n `name` varchar(64) NOT NULL,\n `description` varchar(255) NOT NULL,\n PRIMARY KEY (`name`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nCREATE TABLE `logs` (\n `name` varchar(64) NOT NULL,\n `stamp` datetime NOT NULL,\n `duration` bigint(20) NOT NULL,\n PRIMARY KEY (`name`,`stamp`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nPK\x07\x08\x18\x10\x969M\x01\x00\x00M\x01\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x15\x00	\x0020190619213353.up.sqlUT\x05\x00\x01\x80Cm8ALTER TABLE logs ADD (\n	exit_code  int NOT NULL,\n    output     json NOT NULL\n);\nPK\x07\x08W\xcc\xda\xf9Q\x00\x00\x00Q\x00\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0e\x00	\x00migrations.sqlUT\x05\x00\x01\x80Cm8CREATE TABLE IF NOT EXISTS migrations (\n    project varchar(16) NOT NULL COMMENT 'Microservice or project name',\n    filename varchar(255) NOT NULL COMMENT 'yyyymmddHHMMSS.sql',\n    statement_index int(11) NOT NULL COMMENT 'Statement number from SQL file',\n    status text NOT NULL COMMENT 'ok or full error message',\n\n    PRIMARY KEY (project,filename)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nPK\x07\x08\xca;t\xc5\x89\x01\x00\x00\x89\x01\x00\x00PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\x18\x10\x969M\x01\x00\x00M\x01\x00\x00\x1a\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x00\x00\x00\x0020180704080000.base.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(W\xcc\xda\xf9Q\x00\x00\x00Q\x00\x00\x00\x15\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x9e\x01\x00\x0020190619213353.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\xca;t\xc5\x89\x01\x00\x00\x89\x01\x00\x00\x0e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81;\x02\x00\x00migrations.sqlUT\x05\x00\x01\x80Cm8PK\x05\x06\x00\x00\x00\x00\x03\x00\x03\x00\xe2\x00\x00\x00	\x04\x00\x00\x00\x00"
