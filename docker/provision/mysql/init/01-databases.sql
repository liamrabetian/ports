# create primary and test databases
CREATE DATABASE
    IF NOT EXISTS `ports`;

CREATE DATABASE IF NOT EXISTS `ports_test`;

# create root user and grant rights
CREATE USER
    'root' @'localhost' IDENTIFIED BY 'password';

GRANT ALL PRIVILEGES ON *.* TO 'root'@'%';