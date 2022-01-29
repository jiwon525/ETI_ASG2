package main

//==================== Imports ====================
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Classes struct {
	//first letter must be capitalised.
	ModuleID   string `json:"moduleid"` //get from 3.4``
	ClassID    int    `json:"classid" gorm:"primaryKey"`
	ClassDate  string `json:"classdate"`
	ClassStart string `json:"start_time"`
	ClassEnd   string `json:"end_time"`
	ClassCap   int    `json:"classcap"`
	//TutorID    int    `json:tutorid`1
	TutorName string `json:"tutorname"`
}

const classURL = "http://localhost:9101/api/v1/class"
const key = "2c78afaf-97da-4816-bbee-9ad239abb296"

func CreateClass(c Classes) string {
	// Set up url
	url := classURL + "?key=" + key ///to call. will use to call the apis from other people.
	// Convert to Json
	jsonValue, _ := json.Marshal(c)

	// Post with object
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get fail or success msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 422 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
		}
	}

	response.Body.Close()

	return errMsg
}

func updateClass(ClassID int, c Classes) string {
	// Set up url
	url := classURL + "/" + strconv.Itoa(ClassID) + "?key=" + key ///to call. will use to call the apis from other people.
	// Convert to Json
	jsonValue, _ := json.Marshal(c)

	// Post with object
	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	var errMsg string
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get fail or success msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 422 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
		}
	}
	response.Body.Close()
	return errMsg
}

func Menu() {

	var option string

	// Menu
	fmt.Println("------ Welcome to Class Management ------")
	fmt.Println("[1] Create new class")
	fmt.Println("[2] Update classes")
	fmt.Println("[3] Delete classes")
	fmt.Println("[4] Print all classes")
	fmt.Println("[0] Exit application")
	fmt.Println("------------------------------")
	fmt.Println("Enter your option: ")
	fmt.Scanln(&option)

	// Options
	switch option {
	case "1":
		newClassesMenu()
	case "2":
		updateClassesMenu()
	case "0":

	default:
	}
}
func newClassesMenu() Classes {
	var c Classes

	fmt.Println("CREATING NEW CLASS")
	fmt.Println("Please fill in the following details (b to back).")
	fmt.Println("Module Code: ")
	fmt.Scanln(&c.ModuleID)
	if c.ModuleID == "b" {
		return c
	}
	fmt.Println("Class Day: ")
	fmt.Scanln(&c.ClassDate)
	if c.ClassDate == "b" {
		return c
	}
	fmt.Println("Class Start Time (in HH:MM:SS format): ")
	fmt.Scanln(&c.ClassStart)
	if c.ClassStart == "b" {
		return c
	}
	fmt.Println("Class End Time (in HH:MM:SS format): ")
	fmt.Scanln(&c.ClassEnd)
	fmt.Println("Class Capacity: ")
	fmt.Scanln(&c.ClassCap)
	fmt.Println("Your Name: ")
	fmt.Scanln(&c.TutorName)
	fmt.Println("--------------------")

	// Call api caller to create a new passenger object
	errMsg := CreateClass(c)
	fmt.Println(errMsg)
	if errMsg != "Success" {
		fmt.Println(errMsg)
	}
	return c
}
func updateClassesMenu() Classes {
	var c Classes

	fmt.Println("UPDATING CLASS")
	fmt.Println("Please fill in the following details (b to back).")
	fmt.Println("Class Code: ")
	fmt.Scanln(&c.ClassID)
	if c.ClassID == 0 {
		return c
	}
	fmt.Println("Class Day: ")
	fmt.Scanln(&c.ClassDate)
	if c.ClassDate == "b" {
		return c
	}
	fmt.Println("Class Start Time (in HH:MM:SS format): ")
	fmt.Scanln(&c.ClassStart)
	if c.ClassStart == "b" {
		return c
	}
	fmt.Println("Class End Time (in HH:MM:SS format): ")
	fmt.Scanln(&c.ClassEnd)
	fmt.Println("Class Capacity: ")
	fmt.Scanln(&c.ClassCap)
	fmt.Println("Your Name: ")
	fmt.Scanln(&c.TutorName)
	fmt.Println("--------------------")
	// Call api caller to create a new passenger object
	errMsg := updateClass(c.ClassID, c)
	fmt.Println(errMsg)
	if errMsg != "Success" {
		fmt.Println(errMsg)
	}
	return c
}

func main() {
	//var c Classes

	// While not exit
	for {
		Menu()
	}
}
