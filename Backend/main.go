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
	ClassStart string `json:"start_time"`
	ClassEnd   string `json:"end_time"`
	ClassCap   int    `json:"classcap"`
	//TutorID    int    `json:tutorid`
	TutorName string `json:"tutorname"`
}

var db *sql.DB
var cmap map[int]Classes
var c Classes

//for html

type Class []Classes

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

/*
func allClasses(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/classes_db")
	results, err := db.Query("SELECT * FROM classes_db.Classes")
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()
	var pc []Classes
	for results.Next() {
		// map this type to the record in the table
		var c Classes
		err = results.Scan(&c.ModuleCode, &c.ClassCode, &c.ClassDate, &c.ClassStart, &c.ClassEnd, &c.ClassCap, &c.ClassInfo, &c.ClassRating, &c.TutorName)
		if err != nil {
			panic(err.Error())
		}

		pc = append(pc, c)
		fmt.Println(pc)
		json.NewEncoder(w).Encode(pc)
	}

}
func updateClass(w http.ResponseWriter, r *http.Request) {
	//open db
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/classes_db")
	if err != nil {
		panic(err.Error())
	}
	//defer the db from closing
	defer db.Close()
	fmt.Println("update has been triggered")
	params := mux.Vars(r)
	mod := params["modulecode"]
	sid := params["classcode"]
	id, _ := strconv.Atoi(sid)
	cdate := params["classdate"]
	cstart := params["classstart"]
	cend := params["classend"]
	scap := params["classcap"]
	cap, _ := strconv.Atoi(scap)
	srating := params["classrating"]
	rating, _ := strconv.ParseFloat(srating, 64)
	info := params["classinfo"]
	tname := params["tutorname"]
	fmt.Println(mod, id, cdate, cstart, cend, cap, srating, info, tname)
	//query to create classes
	query := fmt.Sprintf("UPDATE Classes SET ModuleCode='%s',ClassDate='%s',ClassStart='%s',ClassEnd='%s',ClassCap=%d,ClassInfo='%s',ClassRating=%f,TutorName='%s' WHERE ClassCode=%d", mod, cdate, cstart, cend, cap, info, rating, tname, id)
	_, err = db.Query(query)
	//check for error in the query
	if err != nil {
		panic(err.Error())
	}
}

func deleteClass(w http.ResponseWriter, r *http.Request) {
	//open db
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/classes_db")
	if err != nil {
		panic(err.Error())
	}
	//defer the db from closing
	defer db.Close()
	fmt.Println("delete has been triggered")
	//store the content of the request.
	params := mux.Vars(r)
	id := params["classcode"]
	//fmt.Println(id)
	//intid, _ := strconv.Atoi(id)
	//query to delete classes
	query := fmt.Sprintf("DELETE FROM Classes WHERE ClassCode=" + id)
	_, err = db.Query(query)
	//check for error in the query
	if err != nil {
		panic(err.Error())
	}

	//fmt.Fprintf(w, "Class with ClassCode = %s was deleted", params["classcode"])
}*/
func updateClass(w http.ResponseWriter, r *http.Request) {}
func deleteClass(w http.ResponseWriter, r *http.Request) {

}
func classInfo(w http.ResponseWriter, r *http.Request) {

}
func allClasses(w http.ResponseWriter, r *http.Request)    {}
func searchClass(w http.ResponseWriter, r *http.Request)   {}
func classStudents(w http.ResponseWriter, r *http.Request) {}
func main() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/classes_db")

	// Handle error
	if err != nil {
		panic(err.Error())
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/class", createClass).Methods("POST")
	router.HandleFunc("/api/v1/class/{classid}", updateClass).Methods("PUT")
	router.HandleFunc("/api/v1/class/{classid}", deleteClass).Methods("DELETE")
	router.HandleFunc("/api/v1/class/{classid}", classInfo).Methods("GET")                                       //to view class info and ratings
	router.HandleFunc("/api/v1/class", allClasses).Methods("GET")                                                //list all classes
	router.HandleFunc("/api/v1/class?searchKey={searchKey}&filterType={filterType}", searchClass).Methods("GET") //search for classes, filter type to see if filtering by tutor name, class id etc
	router.HandleFunc("/api/v1/class/{classid}", classStudents).Methods("GET")                                   //to view list of students in a class. updated every hour or smt
	fmt.Println("driver microservice api operating on port 9101")
	log.Fatal(http.ListenAndServe(":9101", router))
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "9101"
	}
}

/*
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "9100"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}*/
