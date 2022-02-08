# Table of Contents
1. Introduction
2. Data Structures
3. Microservices
4. Front End
5. Architecture Diagram
6. API Observed


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

![create class](https://user-images.githubusercontent.com/60087854/152891452-74cc4bd1-b84e-40b5-b92a-9539a5dec4c0.png)
![delete class](https://user-images.githubusercontent.com/60087854/152891470-a2b59d34-19cb-4390-b607-a2e5d44709b5.png)
![update class](https://user-images.githubusercontent.com/60087854/152895345-85f3f6d4-6c93-4e0c-b41c-3be60b1065ee.png)
![searching classes](https://user-images.githubusercontent.com/60087854/152895371-39effcf3-239d-4e35-9707-2218f9274706.png)
![list of classes](https://user-images.githubusercontent.com/60087854/152895397-b6bc1ba0-31d6-4106-8570-3f4d3e1176ea.png)
![List of Students](https://user-images.githubusercontent.com/60087854/152891438-d5eba434-c044-4e5b-92b4-3279e729dbaa.png)

## 5. Architecture Diagram
![diagram](https://user-images.githubusercontent.com/60087854/152892550-046c6331-1c45-4ac4-8b22-047ee2a2bec7.png)
- [x] Calling in data from Ratings to get the average rating value to print </br>
- [x] Calling in data from Module to get the module synopsis with module code identifier </br>
- [ ] Calling in Authentication to get Tutor ID for creation of class</br>
--------------------------------------------------------------------------
Bid Service calles my apis to get list of classes </br>
Timetable service calls my apis to get class details

## 6. API Observed

| URL | Method  | Description  |
| :--------: | :---: | :----------: |
| /api/v1:9051/api/v1/allocations/class/{classid} | GET | Getting a list of students in a certain class |
| "http://localhost:9141/api/v1/module/{modulecode} | GET | Get module synopsis |
| http://10.31.11.12:9042/api/rating/class/{classid} | GET | Get all list of ratings for a certain class |
