package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Classes struct {
	//first letter must be capitalised.
	ModuleID   string `json:"moduleid"` //get from 3.4``
	ClassID    int    `json:"classid" gorm:"primaryKey"`
	ClassDate  string `json:"classdate"`
	ClassStart string `json:"classstart"`
	ClassEnd   string `json:"classend"`
	ClassCap   int    `json:"classcap"`
	TutorName  string `json:"tutorname"`
	/*TutorID    int    `json:"tutorid"`
	Rating float32
	ClassInfo string */
}

//call other peoples api for TutorID TutorName rating classinfo

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

	//PassengerID is auto incremented
	query := fmt.Sprintf("INSERT INTO Classes (ModuleID,ClassDate,ClassStart,ClassEnd,ClassCap,TutorName) VALUES('%s','%s','%s','%s',%d,'%s')", c.ModuleID, c.ClassDate, c.ClassStart, c.ClassEnd, c.ClassCap, c.TutorName)
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
			if c.ModuleID == "" || c.ClassDate == "" || c.ClassStart == "" || c.ClassEnd == "" || c.ClassCap == 0 || c.TutorName == "" {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply all neccessary information to create new class "))
			} else {
				//all necessary info inside
				// Run db create class func
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
	query := fmt.Sprintf("UPDATE Classes SET ModuleID='%s',ClassDate='%s',ClassStart='%s',ClassEnd='%s',ClassCap=%d,TutorName='%s' WHERE ClassID=%d", c.ModuleID, c.ClassDate, c.ClassStart, c.ClassEnd, c.ClassCap, c.TutorName, cid)
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
			if c.ClassID == 0 || c.ModuleID == "" || c.ClassDate == "" || c.ClassStart == "" || c.ClassEnd == "" || c.ClassCap == 0 || c.TutorName == "" {
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
func getclassDB(db *sql.DB) ([]Classes, string) {
	query := fmt.Sprintf("SELECT * FROM classes_db.classes")
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var errMsg string
	var cc []Classes
	for results.Next() {
		var c Classes
		err = results.Scan(&c.ClassID, &c.ModuleID, &c.ClassDate, &c.ClassStart, &c.ClassEnd, &c.ClassCap, &c.TutorName)
		if err != nil {
			errMsg = "Classes do not exist"
		}
		cc = append(cc, c)
		//fmt.Println(cc)
	}
	return cc, errMsg
}

//api to print all classes
func allClasses(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		c, errMsg := getclassDB(db)

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
func searchClass(w http.ResponseWriter, r *http.Request) {

}
func classStudents(w http.ResponseWriter, r *http.Request) {}
func classInfo(w http.ResponseWriter, r *http.Request) {

}
func main() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/classes_db")

	// Handle error
	if err != nil {
		panic(err.Error())
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/class", createClass).Methods("POST")
	//router.HandleFunc("/api/v1/class/{tutorid}",tutorsearch).Methods("GET")
	router.HandleFunc("/api/v1/class/{classid}", updateClass).Methods("PUT")
	router.HandleFunc("/api/v1/class/{classid}", deleteClass).Methods("DELETE")
	router.HandleFunc("/api/v1/class/{classid}", classInfo).Methods("GET")               //to view class info and ratings
	router.HandleFunc("/api/v1/class", allClasses).Methods("GET")                        //list all classes
	router.HandleFunc("/api/v1/class?searchKey={searchKey}", searchClass).Methods("GET") //search for classes, filter type to see if filtering by tutor name, class id etc
	router.HandleFunc("/api/v1/class/{classid}", classStudents).Methods("GET")           //to view list of students in a class. updated every hour or smt
	fmt.Println("driver microservice api operating on port 9101")
	log.Fatal(http.ListenAndServe(":9101", router))
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "9101"
	}
}
