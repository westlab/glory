DROP DATABASE IF EXISTS glory_test;
CREATE DATABASE IF NOT EXISTS glory_test;
USE glory_test;

DROP TABLE IF EXISTS thesis_history;
DROP TABLE IF EXISTS author;

CREATE TABLE author (
  author_id INT(11) AUTO_INCREMENT NOT NULL PRIMARY KEY,
  name CHAR(32) NOT NULL,
  working_group CHAR(32) NOT NULL
) ENGINE=INNODB DEFAULT CHARSET=utf8;

CREATE TABLE thesis_history (
  id INT(11) AUTO_INCREMENT NOT NULL PRIMARY KEY,
  author_id INT(11) NOT NULL,
  char_count INT(7) NOT NULL,
  last_mod DATETIME(3) NOT NULL,
  fetch_time DATETIME(3) NOT NULL,
  CONSTRAINT fk_author_id FOREIGN KEY (author_id) REFERENCES author (author_id) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=INNODB DEFAULT CHARSET=utf8;

GRANT ALL PRIVILEGES ON *.* TO 'westlab'@'%';
