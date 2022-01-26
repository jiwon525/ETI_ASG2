create table Classes(ClassID int NOT NULL AUTO_INCREMENT, 
    ModuleID VARCHAR(5) NOT NULL,
	ClassDate varchar(10) NOT NULL,  
	ClassStart  time NOT NULL, 
	ClassEnd time NOT NULL,
    ClassCap int NOT NULL,  
	TutorName VARCHAR(30) NOT NULL,
    PRIMARY KEY (ClassID));