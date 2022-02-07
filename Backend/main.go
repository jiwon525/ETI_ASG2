package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Classes struct {
	//first letter must be capitalised.
	ClassID    int     `json:"classid"`
	ModuleCode string  `json:"modulecode"`
	ClassDate  string  `json:"classdate"`
	ClassStart string  `json:"classstart"`
	ClassEnd   string  `json:"classend"`
	ClassCap   int     `json:"classcap"`
	TutorName  string  `json:"tutorname"`
	TutorID    int     `json:"tutorid"`
	Rating     float64 `json:"rating"`
	ClassInfo  string  `json:"classinfo"`
}

//to accept incoming data from external api that i call--------------
type Modules struct {
	ModuleID          int    `gorm:"primaryKey"`
	ModuleCode        string `json:"modulecode"`
	ModuleName        string `json:"modulename"`
	Synopis           string `json:"synopis"`
	LearningObjective string `json:"learningobjective"`
	Deleted           gorm.DeletedAt
}
type Rating struct {
	RatingID          int
	CreatorID         int
	CreatorType       string
	TargetID          int
	TargetType        string
	RatingScore       int
	Anonymous         int
	DateTimePublished string
	CreatorName       string
	TargetName        string
}
type Student struct {
	StudentID int
	Semester  string
	ClassID   int
}
type Tutor struct {
	Deleted     gorm.DeletedAt
	TutorID     int    `json:"tutor_id" gorm:"primaryKey"`
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Description string `json:"descriptions" validate:"required"`
}

//-----------------------------------------------------------------------

//----------function for api calling from other packages----------
//calling timetable module to get a list of students
func classStudents(cid int) (string, []Student) {
	url := "http://10.31.11.12:9051/api/v1/allocations/class/" + strconv.Itoa(cid)
	// Create request
	response, err := http.Get(url)
	var s []Student
	var errMsg string
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		// Fetch Request
		data, _ := ioutil.ReadAll(response.Body)
		// Get fail or success msg
		if response.StatusCode == 422 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
			json.Unmarshal([]byte(data), &s) //getting just the student id. need to get name from 3.5
		}
	}
	response.Body.Close()
	return errMsg, s
}

//api caller for class rating from package 3.9
func callClassInfo(modid string) string {
	//set up url
	url := "http://localhost:9141/api/v1/module/" + modid
	// Create request
	response, err := http.Get(url)
	var mods Modules
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		// Fetch Request
		data, _ := ioutil.ReadAll(response.Body)
		// Get fail or success msg
		if response.StatusCode == 422 {
			fmt.Println(string(data))
		} else {
			json.Unmarshal([]byte(data), &mods)
		}
	}
	response.Body.Close()
	return mods.Synopis
}

//api caller for module synopsis from package 3.4 Management of Modules
func callClassRating(cid int) (string, []Rating) {
	//set up url
	url := "http://10.31.11.12:9042/api/rating/class/" + strconv.Itoa(cid)
	// Create request
	response, err := http.Get(url)
	var r []Rating
	var errMsg string
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		// Fetch Request
		data, _ := ioutil.ReadAll(response.Body)
		// Get fail or success msg
		if response.StatusCode == 422 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
			json.Unmarshal([]byte(data), &r)
		}
	}
	response.Body.Close()
	return errMsg, r
}

//to collate an average score from the rating data collected
func ratingaverage(r []Rating) float64 {
	sum := 0
	score := 0
	for _, rr := range r {
		score = score + rr.RatingScore
		sum++
	}

	avg := float64(score) / float64(sum)
	return avg

}

//-----------------------------------------------------------
var db *sql.DB

/*for html
const classURL = "http://10.31.11.12:9101/api/v1/class"
const key = "2c78afaf-97da-4816-bbee-9ad239abb296"*/

//==================== Auxiliary Functions ====================
func validKey(r *http.Request) bool {
	v := r.URL.Query()
	if key, ok := v["key"]; ok {
		if key[0] == "2c78afaf-97da-4816-bbee-9ad239abb296" {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

//=====================My REST APIs=======================
//function to insert new classes into db (POST)
func createClassDB(db *sql.DB, c Classes) {

	//primary key class id is auto incremented.
	//query to insert into db
	query := fmt.Sprintf("INSERT INTO Classes (ModuleCode,ClassDate,ClassStart,ClassEnd,ClassCap,TutorName,TutorID,Rating,ClassInfo) VALUES('%s','%s','%s','%s',%d,'%s',%d,%f,'%s')", c.ModuleCode, c.ClassDate, c.ClassStart, c.ClassEnd, c.ClassCap, c.TutorName, c.TutorID, c.Rating, c.ClassInfo)
	_, err := db.Query(query)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err.Error())
	}
	fmt.Println("Successfully added into the DB")
}

//api to create new class and call the createClassDB function
func createClass(w http.ResponseWriter, r *http.Request) {
	// Valid key for API check
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}
	if r.Header.Get("Content-type") == "application/json" {
		reqBody, err := ioutil.ReadAll(r.Body)

		if err == nil { // If no error
			var c Classes
			// Map json to variable Classes c
			json.Unmarshal([]byte(reqBody), &c)

			// Check if all non-null information exist
			if c.ModuleCode == "" || c.ClassDate == "" || c.ClassStart == "" || c.ClassEnd == "" || c.ClassCap == 0 || c.TutorName == "" || c.TutorID == 0 {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply all neccessary information to create new class "))
			} else {
				//all necessary info inside
				createClassDB(db, c)
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("201 - Class created."))

			}
		} else { //Incorrect format
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply class information in JSON format"))
		}
	}
}

//function to update class details in the db
func updateClassDB(db *sql.DB, cid int, c Classes) {
	query := fmt.Sprintf("UPDATE Classes SET ModuleCode='%s',ClassDate='%s',ClassStart='%s',ClassEnd='%s',ClassCap=%d,TutorName='%s',TutorID=%d WHERE ClassID=%d", c.ModuleCode, c.ClassDate, c.ClassStart, c.ClassEnd, c.ClassCap, c.TutorName, c.TutorID, cid)
	_, err := db.Query(query)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err.Error())
	}
	fmt.Println("Successfully updated the DB")
}

//api to update class details and call the updateClassDB
func updateClass(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}
	params := mux.Vars(r)
	var cid int
	fmt.Sscan(params["classid"], &cid)

	if r.Header.Get("Content-type") == "application/json" {
		reqBody, err := ioutil.ReadAll(r.Body)

		if err == nil { // If no error
			var c Classes
			// Map json to variable Classes c
			json.Unmarshal([]byte(reqBody), &c)
			// Check if all non-null information exist
			if c.ClassID == 0 || c.ModuleCode == "" || c.ClassDate == "" || c.ClassStart == "" || c.ClassEnd == "" || c.ClassCap == 0 || c.TutorName == "" || c.TutorID == 0 {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply all neccessary information to update existing class "))
			} else {
				updateClassDB(db, cid, c)
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("204 - Resource updated successfully"))
			}
		} else { //Incorrect format
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply class information in JSON format"))
		}
	}
}

//function to delete class from the database
func deleteClassDB(db *sql.DB, cid int) string {
	query := fmt.Sprintf("DELETE FROM Classes WHERE ClassID=%d", cid)
	_, err := db.Query(query)
	errMsg := "Success"
	if err != nil {
		errMsg = "Class does not exist"
	}
	fmt.Println("Successfully deleted item from DB")
	return errMsg
}

//api to delete classes
func deleteClass(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		params := mux.Vars(r)
		var cid int
		fmt.Sscan(params["classid"], &cid)
		errMsg := deleteClassDB(db, cid)
		if errMsg == "Success" {
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Class deleted: " +
				params["classid"]))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No class found"))
		}

	}

}

func deleteModDB(db *sql.DB, modid string) string {
	query := fmt.Sprintf("DELETE FROM Classes WHERE ModuleCode='%s'", modid)
	_, err := db.Query(query)
	errMsg := "Success"
	if err != nil {
		errMsg = "Class does not exist"
	}
	fmt.Println("Successfully deleted item from DB")
	return errMsg
}

func deleteMod(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		params := mux.Vars(r)
		var modid string
		fmt.Sscan(params["ModuleCode"], &modid)
		errMsg := deleteModDB(db, modid)
		if errMsg == "Success" {
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Classes related to Module Code deleted: " +
				params["ModuleCode"]))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No class found"))
		}

	}

}

//function to connect with database to print it out
func getClassDB(db *sql.DB) ([]Classes, string) {
	results, err := db.Query("SELECT * FROM classes_db.classes")
	if err != nil {
		panic(err.Error())
	}
	var errMsg string
	var cc []Classes
	for results.Next() {
		var c Classes
		err = results.Scan(&c.ClassID, &c.ModuleCode, &c.ClassDate, &c.ClassStart, &c.ClassEnd, &c.ClassCap, &c.TutorName, &c.TutorID)
		if err != nil {
			errMsg = "Classes do not exist"
		}
		cinfo := callClassInfo(c.ModuleCode)
		errMsg3, r := callClassRating(c.ClassID)
		errMsg = errMsg3
		rating := ratingaverage(r)
		c.ClassInfo = cinfo
		c.Rating = rating
		cc = append(cc, c)
	}
	return cc, errMsg
}

//api to print all classes
func allClasses(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		c, errMsg := getClassDB(db)

		switch errMsg {
		case "Classes do not exist":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No class found"))
		default:
			//call api
			json.NewEncoder(w).Encode(c)
		}

	}

}
func searchClassDB(db *sql.DB, mid string) ([]Classes, string) {
	query := fmt.Sprintf("SELECT * FROM classes_db.classes WHERE ModuleCode='%s'", mid)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var errMsg string
	var cc []Classes
	for results.Next() {
		var c Classes
		err = results.Scan(&c.ClassID, &c.ModuleCode, &c.ClassDate, &c.ClassStart, &c.ClassEnd, &c.ClassCap, &c.TutorName, &c.TutorID)
		fmt.Println(c)
		if err != nil {
			errMsg = "Classes do not exist"
		}
		cinfo := callClassInfo(c.ModuleCode)
		errMsg3, r := callClassRating(c.ClassID)
		errMsg = errMsg3
		rating := ratingaverage(r)
		c.ClassInfo = cinfo
		c.Rating = rating
		cc = append(cc, c)
		//fmt.Println(cc)
	}
	return cc, errMsg
}
func searchClass(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var mid string
		// Get query string parameters of student ID and semester start date
		queryString := r.URL.Query()

		fmt.Sscan(queryString["ModuleCode"][0], &mid)

		c, errMsg := searchClassDB(db, mid)

		switch errMsg {
		case "Classes do not exist":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No class found"))
		default:
			//call api
			json.NewEncoder(w).Encode(c)
		}

	}
}

func getStudentList(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var cid int
		// Get query string parameters of student ID and semester start date
		queryString := r.URL.Query()

		fmt.Sscan(queryString["classid"][0], &cid)

		errMsg, student := classStudents(cid)
		switch errMsg {
		case "Students do not exist":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No class found"))
		default:
			//call api
			json.NewEncoder(w).Encode(student)
		}

	}
}

// Help function that calls appropriate function in accordance to parameters in the query string
func GetClassQuery(w http.ResponseWriter, r *http.Request) {
	// Get query string parameters
	queryString := r.URL.Query()
	_, searchclass := queryString["ModuleCode"]
	if searchclass {
		searchClass(w, r)
		return
	} else {
		allClasses(w, r)
		return
	}
}

// Help function that calls appropriate function in accordance to parameters in the query string
func DeleteClassQuery(w http.ResponseWriter, r *http.Request) {
	// Get query string parameters
	queryString := r.URL.Query()
	_, deletemod := queryString["ModuleCode"]
	_, deleteclass := queryString["classid"]
	if deleteclass {
		deleteClass(w, r)
		return
	} else if deletemod {
		deleteMod(w, r)
		return
	}
}

func main() {
	var err error
	db, err = sql.Open("mysql", "root:classdatabase@tcp(classdatabase:3306)/classes_db")

	// Handle error
	if err != nil {
		panic(err.Error())
	}

	router := mux.NewRouter()
	// This is to allow the headers, origins and methods all to access CORS resource sharing
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})

	router.HandleFunc("/api/v1/class", createClass).Methods("POST")
	router.HandleFunc("/api/v1/class/{classid}", updateClass).Methods("PUT")
	router.HandleFunc("/api/v1/class/", DeleteClassQuery).Methods("DELETE")
	router.HandleFunc("/api/v1/class/{classid}", getStudentList).Methods("GET")
	router.HandleFunc("/api/v1/class", GetClassQuery).Methods("GET")
	fmt.Println("driver microservice api operating on port 9101")
	log.Fatal(http.ListenAndServe(":9101", handlers.CORS(headers, origins, methods)(router)))
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "9101"
	}
}
