DROP DATABASE IF EXISTS `classes_db`;
CREATE DATABASE `classes_db`;
USE `classes_db`;

create table `Classes`(`ClassID` int NOT NULL AUTO_INCREMENT,
	`ModuleCode` VARCHAR(5) NOT NULL,
	`ClassDate` varchar(10),
	`ClassStart`  varchar(4),
	`ClassEnd` varchar(4),
    `ClassCap` int,
	`TutorName` VARCHAR(30),
    `TutorID` int,
    PRIMARY KEY (ClassID)
);