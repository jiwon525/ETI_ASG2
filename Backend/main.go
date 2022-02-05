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

//need to edit localhost to the front end ip address
type Classes struct {
	//first letter must be capitalised.
	ClassID    int     `json:"classid"`  // gorm:"primaryKey"
	ModuleID   string  `json:"moduleid"` //my module ID = module code from package 3.4
	ClassDate  string  `json:"classdate"`
	ClassStart string  `json:"classstart"`
	ClassEnd   string  `json:"classend"`
	ClassCap   int     `json:"classcap"`
	TutorFName string  `json:"tutorfname"` //need to edit main.go
	TutorLName string  `json:"tutorlname"` //need to edit main.go
	TutorID    int     `json:"tutorid"`    //need to edit main.go
	Rating     float64 `json:"rating"`     //need to edit main.go, get from 3.9
	ClassInfo  string  `json:"classinfo"`  //get from 3.4
}

//to accept incoming data by calling api from package 3.4
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
	Sname     string
	ClassID   int
}

//call other peoples api for TutorID TutorName rating classinfo
//function for api calling from other packages
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
			json.Unmarshal([]byte(data), &s)
		}
	}
	response.Body.Close()
	return errMsg, s
}

//api caller for class rating from package 3.9
func callClassInfo(modid string) (string, string) {
	//set up url
	url := "http://localhost:9141/api/v1/module/" + modid
	// Create request
	response, err := http.Get(url)
	var mods Modules
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
			json.Unmarshal([]byte(data), &mods)
		}
	}
	response.Body.Close()
	return errMsg, mods.Synopis
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

var db *sql.DB
var c Classes

//for html

const classURL = "http://localhost:9101/api/v1/class"
const key = "2c78afaf-97da-4816-bbee-9ad239abb296"

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

//function to insert new classes in the db
func createClassDB(db *sql.DB, c Classes) {

	//primary key class id is auto incremented.
	//query to insert into db
	query := fmt.Sprintf("INSERT INTO Classes (ModuleID,ClassDate,ClassStart,ClassEnd,ClassCap,TutorFName,TutorLName,TutorID,Rating,ClassInfo) VALUES('%s','%s','%s','%s',%d,'%s','%s',%d,%f,'%s')", c.ModuleID, c.ClassDate, c.ClassStart, c.ClassEnd, c.ClassCap, c.TutorFName, c.TutorLName, c.TutorID, c.Rating, c.ClassInfo)
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
			if c.ModuleID == "" || c.ClassDate == "" || c.ClassStart == "" || c.ClassEnd == "" || c.ClassCap == 0 || c.TutorFName == "" || c.TutorLName == "" || c.TutorID == 0 {
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
	query := fmt.Sprintf("UPDATE Classes SET ModuleID='%s',ClassDate='%s',ClassStart='%s',ClassEnd='%s',ClassCap=%d,TutorFName='%s',TutorLName='%s',TutorID=%d WHERE ClassID=%d", c.ModuleID, c.ClassDate, c.ClassStart, c.ClassEnd, c.ClassCap, c.TutorFName, c.TutorLName, c.TutorID, cid)
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
			if c.ClassID == 0 || c.ModuleID == "" || c.ClassDate == "" || c.ClassStart == "" || c.ClassEnd == "" || c.ClassCap == 0 || c.TutorFName == "" || c.TutorLName == "" || c.TutorID == 0 {
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
			w.Write([]byte("202 - Course deleted: " +
				params["courseid"]))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No class found"))
		}

	}

}

//function to connect with database to print it out
func getClassDB(db *sql.DB) ([]Classes, string) {
	query := fmt.Sprintf("SELECT * FROM classes_db.classes")
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var errMsg string
	var cc []Classes
	for results.Next() {
		var c Classes
		err = results.Scan(&c.ClassID, &c.ModuleID, &c.ClassDate, &c.ClassStart, &c.ClassEnd, &c.ClassCap, &c.TutorFName, &c.TutorLName, &c.TutorID)
		if err != nil {
			errMsg = "Classes do not exist"
		}
		errMsg2, cinfo := callClassInfo(c.ModuleID)
		errMsg = errMsg2
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
	query := fmt.Sprintf("SELECT * FROM classes_db.classes WHERE ModuleID='%s'", mid)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var errMsg string
	var cc []Classes
	for results.Next() {
		var c Classes
		err = results.Scan(&c.ClassID, &c.ModuleID, &c.ClassDate, &c.ClassStart, &c.ClassEnd, &c.ClassCap, &c.TutorFName, &c.TutorLName, &c.TutorID)
		fmt.Println(c)
		if err != nil {
			errMsg = "Classes do not exist"
		}
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

		fmt.Sscan(queryString["ModuleID"][0], &mid)

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

// Help function that calls appropriate function in accordance to parameters in the query string
func GetClassQuery(w http.ResponseWriter, r *http.Request) {
	// Get query string parameters
	queryString := r.URL.Query()
	_, searchclass := queryString["ModuleID"]
	if searchclass {
		searchClass(w, r)
		return
	} else {
		allClasses(w, r)
		return
	}
}
func main() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/classes_db")

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
	router.HandleFunc("/api/v1/class/{classid}", deleteClass).Methods("DELETE")
	router.HandleFunc("/api/v1/class/{classid}", searchClass).Methods("GET")
	router.HandleFunc("/api/v1/class", GetClassQuery).Methods("GET")
	fmt.Println("driver microservice api operating on port 9101")
	log.Fatal(http.ListenAndServe(":9101", handlers.CORS(headers, origins, methods)(router)))
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "9101"
	}
}
