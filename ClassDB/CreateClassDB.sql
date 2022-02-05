create table Classes(ClassID int NOT NULL AUTO_INCREMENT, 
	ModuleID VARCHAR(5) NOT NULL,
	ClassDate varchar(10),  
	ClassStart  varchar(4), 
	ClassEnd varchar(4),
    ClassCap int,  
	TutorFName VARCHAR(30),
    TutorLName VARCHAR(30),
    TutorID int,
    PRIMARY KEY (ClassID));