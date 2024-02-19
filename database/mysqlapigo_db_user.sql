-- DCL (Data Control Language)
-- create user 
CREATE USER 'mysql_apigo_user'@'localhost' IDENTIFIED BY 'MySql_ApiGo#97';

-- grant privileges
GRANT ALL PRIVILEGES ON mysqlapigo.* TO 'mysql_apigo_user'@'localhost';

-- flush privileges
FLUSH PRIVILEGES;


