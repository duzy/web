CREATE DATABASE dusell
CREATE DATABASE test
CREATE USER 'test'@'localhost' IDENTIFIED BY 'abc'
GRANT ALL PRIVILEGES ON dusell.* TO 'test'@'localhost'
GRANT ALL ON test.* TO 'test'@'localhost'

-- SET PASSWORD FOR 'test'@'localhost' = PASSWORD('abc')
-- DROP USER 'test'@'localhost'
