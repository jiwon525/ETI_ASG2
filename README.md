# Table of Contents
1. Introduction


## 1. Introduction
ETI Assignment 2. <br></br>
The package that I am working on is 3.8, Management of Classes with specific requirements as shown below.
```
Requirements:
3.8.1.	Create, update, delete classes (must incl. Class Code, Class Schedule, Class Capacity)
3.8.2.	View class info and ratings
3.8.3.	List all classes
3.8.4.	Search for classes
3.8.5.	List all students in a class
```
Only tutors who have logged in would be able to create update delete classes, and the other functions will be open to students.

## 2. Data Structures
Data in package 3.8 has been laid out in this way.</br> *Ratings and ClassInfo are not saved inside the database, they are called from other packages and are only used during my 'List all classes' function.*

| Field Name | Type  | Description  |
| :--------: | :---: | :----------: |
| ClassID | int | The unique primary key that is auto incremented and used to identify classes |
| ModuleCode | string | The acronym for classes which helps to identify which module this class teaches. |
| ClassDate | string | The day of the week that the lessons occurs on. (E.g. Monday) |
| ClassStart | string | Data written in 24-HR clock format to show when the class starts |
| ClassEnd | string | Data written in 24-HR clock format to show when the class ends |
| ClassCap | int | Total student capacity of the class |
| TutorName | string | Tutor's name |
| TutorID | int | Unique ID to identify the tutor |
| Rating | float64 | The rating of the class |
| ClassInfo | string | The module synopsis as extra info related to class |

Json Version:
```
{
	classid:1
	modulecode:"DF"
	classdate:"Tuesday"
	classstart:"0900"
	classend:"1300"
	classcap:30
	tutorname:"Wong Liew"
	tutorid:1
	rating:3.4
	classinfo: "DF"
}
```
MYSQL Table:</br>
{create table Classes(
	ClassID int NOT NULL AUTO_INCREMENT, 
	ModuleID VARCHAR(5) NOT NULL,
	ClassDate varchar(10),  
	ClassStart  varchar(4), 
	ClassEnd varchar(4),
    ClassCap int,  
	TutorFName VARCHAR(30),
    TutorID int,
    PRIMARY KEY (ClassID));
}
## 3. Microservices 
Base URL for class: http://10.31.11.12:9101

| URL | Method  | Description  |
| :--------: | :---: | :----------: |
| /api/v1/class | POST | Create new class, all except class ID needs to be supplied |
| /api/v1/class?classid={classid} | DELETE | Delete class by Class ID |
| /api/v1/class?modulecode={modulecode} | DELETE | Delete all classes with certain Module Code |
| /api/v1/class/{classid} | PUT | Update class by Class ID |
| /api/v1/class/{classid} | GET | Get list of students in a class by its Class ID |
| /api/v1/class | GET | Get all list of classes currently in the database |
| /api/v1/class?ModuleCode={modulecode} | GET | Get list of classes teaching a certain module |

**All the above APIs will require the usage of an authentication key through a query string: ?key = "2c78afaf-97da-4816-bbee-9ad239abb298"**

## 4. Front End


![image](/Images/List of Students.png)

### functions for now
<ul><li>delete class:
curl -X DELETE http://10.31.11.12:9101/api/v1/class/7?key=2c78afaf-97da-4816-bbee-9ad239abb296</li>
<li>create class:
curl -X POST -H "Content-Type:application/json" -d "{"moduleid":"DF","classdate":"Monday","classstart":"11:00:00","classend":"13:00:00","classcap":30,"tutorname":"James_Lee"}" "http://localhost:9101/api/v1/class?key=2c78afaf-97da-4816-bbee-9ad239abb296"</li>
<li>update class(with class ID=1):
curl -H "Content-Type:application/json" -X PUT http://localhost:9101/api/v1/class/1?key=2c78afaf-97da-4816-bbee-9ad239abb296 -d "{"moduleid":"DF","classdate":"Monday","classstart":"11:00:00","classend":"13:00:00","classcap":30,"tutorname":"James_Lee"}"</li>
</ul>

