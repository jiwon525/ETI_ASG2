package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Classes struct {
	//first letter must be capitalised.
	ModuleCode  string  `json:"modulecode"` //get from 3.4``
	ClassCode   int     `json:"classcode" gorm:"primaryKey"`
	ClassDate   string  `json:"classdate"`
	ClassStart  string  `json:"start_time"`
	ClassEnd    string  `json:"end_time"`
	ClassCap    int     `json:"classcap"`
	ClassInfo   string  `json:"classinfo"`
	ClassRating float32 `json:"classrating"` //get from 3.9?
	TutorName   string  `json:"tutorname"`   //get from 3.3?
}

var db *sql.DB
var cmap map[int]Classes

//for html

type Class []Classes

func createClass(w http.ResponseWriter, r *http.Request) {
	//open db
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/classes_db")
	if err != nil {
		panic(err.Error())
	}
	//defer the db from closing
	defer db.Close()
	fmt.Println("create has been triggered")
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
	query := fmt.Sprintf("INSERT INTO Classes VALUES('%s',%d,'%s','%s','%s',%d,'%s',%f,'%s')", mod, id, cdate, cstart, cend, cap, info, rating, tname)
	_, err = db.Query(query)
	//check for error in the query
	if err != nil {
		panic(err.Error())
	}
}
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
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/createclass/{modulecode}/{classcode}/{classdate}/{classstart}/{classend}/{classcap}/{classinfo}/{classrating}/{tutorname}", createClass).Methods("POST")
	router.HandleFunc("/api/editclass/{modulecode}/{classcode}/{classdate}/{classstart}/{classend}/{classcap}/{classinfo}/{classrating}/{tutorname}", updateClass).Methods("PUT")
	router.HandleFunc("/api/deleteclass/{classcode}", deleteClass).Methods("DELETE")
	router.HandleFunc("/api/allclass", allClasses).Methods("GET")
	fmt.Println("driver microservice api operating on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

	// http.HandleFunc("/api/deleteclass", deleteClass)
	// http.HandleFunc("/", home)
	// fmt.Println("server open")
	// http.ListenAndServe(":8080", nil)

}
