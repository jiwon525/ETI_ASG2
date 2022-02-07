# Table of Contents
1. Introduction


### 1. Introduction
Code for ETI Assignment 2. 
Package I am working on : 3.8 Management of Classes
'''Requirements:
3.8.1.	Create, update, delete classes. Info includes
3.8.1.1.	Class code
3.8.1.2.	Schedule of the class
3.8.1.3.	Capacity of class
3.8.2.	View class info and ratings can do this  in listing classes??
3.8.3.	List all classes
3.8.4.	Search for classes
3.8.5.	List all students in a class'''



### functions for now
<ul><li>delete class:
curl -X DELETE http://10.31.11.12:9101/api/v1/class/7?key=2c78afaf-97da-4816-bbee-9ad239abb296</li>
<li>create class:
curl -X POST -H "Content-Type:application/json" -d "{"moduleid":"DF","classdate":"Monday","classstart":"11:00:00","classend":"13:00:00","classcap":30,"tutorname":"James_Lee"}" "http://localhost:9101/api/v1/class?key=2c78afaf-97da-4816-bbee-9ad239abb296"</li>
<li>update class(with class ID=1):
curl -H "Content-Type:application/json" -X PUT http://localhost:9101/api/v1/class/1?key=2c78afaf-97da-4816-bbee-9ad239abb296 -d "{"moduleid":"DF","classdate":"Monday","classstart":"11:00:00","classend":"13:00:00","classcap":30,"tutorname":"James_Lee"}"</li>
</ul>

### final struct format would be:
	ClassID    int     `json:"classid"`  
	ModuleID   string  `json:"moduleid"` 
	ClassDate  string  `json:"classdate"`
	ClassStart string  `json:"classstart"`
	ClassEnd   string  `json:"classend"`
	ClassCap   int     `json:"classcap"`
	TutorFName string  `json:"tutorfname"` 
	TutorLName string  `json:"tutorlname"` 
	TutorID    int     `json:"tutorid"`    
	Rating     float64 `json:"rating"`     
	ClassInfo  string  `json:"classinfo"`  
}

### Database Values
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