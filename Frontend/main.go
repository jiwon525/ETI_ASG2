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
	ClassID    int     `json:"classid"`  // gorm:"primaryKey"
	ModuleID   string  `json:"moduleid"` //my module ID = module code from package 3.4
	ClassDate  string  `json:"classdate"`
	ClassStart string  `json:"classstart"`
	ClassEnd   string  `json:"classend"`
	ClassCap   int     `json:"classcap"`
	TutorFName string  `json:"tutorfname"` //need to edit main.go
	TutorLName string  `json:"tutorlname"` //need to edit main.go
	TutorID    int     `json:"tutorid"`    //need to edit main.go
	Rating     float32 `json:"rating"`     //need to edit main.go, get from 3.9
	ClassInfo  string  `json:"classinfo"`  //get from 3.4
}

const classURL = "http://localhost:9101/api/v1/class"
const key = "2c78afaf-97da-4816-bbee-9ad239abb296"

func printSlice(s []Classes) {
	for _, el := range s {
		fmt.Println(el)
	}

}
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
	//putting new data
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
func deleteClass(ClassID int) string {
	//set up url
	url := classURL + "/" + strconv.Itoa(ClassID) + "?key" + key
	// Request (DELETE "http://localhost:9101/api/v1/class/x?key/<key value>")

	// Create request
	request, err := http.NewRequest("DELETE", url, nil)
	// Create client
	client := &http.Client{}
	response, err := client.Do(request)
	var errMsg string
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		// Fetch Request
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

func newClassesMenu() {
	var c Classes

	fmt.Println("CREATING NEW CLASS")
	fmt.Println("Please fill in the following details.")
	fmt.Println("Module Code: ")
	fmt.Scanln(&c.ModuleID)
	fmt.Println("Class Day: ")
	fmt.Scanln(&c.ClassDate)
	fmt.Println("Class Start Time (in HH:MM:SS format): ")
	fmt.Scanln(&c.ClassStart)
	fmt.Println("Class End Time (in HH:MM:SS format): ")
	fmt.Scanln(&c.ClassEnd)
	fmt.Println("Class Capacity: ")
	fmt.Scanln(&c.ClassCap)
	fmt.Println("Your FirstName: ")
	fmt.Scanln(&c.TutorFName)
	fmt.Println("Your Name: ")
	fmt.Scanln(&c.TutorFName)
	fmt.Println("--------------------")

	// Call api caller to create a new passenger object
	errMsg := CreateClass(c)
	fmt.Println(errMsg)
	if errMsg != "Success" {
		fmt.Println(errMsg)
	}
}
func updateClassesMenu() {
	var c Classes

	fmt.Println("UPDATING CLASS")
	fmt.Println("Please fill in the following details.")
	fmt.Println("Class Code: ")
	fmt.Scanln(&c.ClassID)
	fmt.Println("Module Code: ")
	fmt.Scanln(&c.ModuleID)
	fmt.Println("Class Day: ")
	fmt.Scanln(&c.ClassDate)
	fmt.Println("Class Start Time (in HH:MM:SS format): ")
	fmt.Scanln(&c.ClassStart)
	fmt.Println("Class End Time (in HH:MM:SS format): ")
	fmt.Scanln(&c.ClassEnd)
	fmt.Println("Class Capacity: ")
	fmt.Scanln(&c.ClassCap)
	fmt.Println("Your First Name: ")
	fmt.Scanln(&c.TutorFName)
	fmt.Println("Your Last Name: ")
	fmt.Scanln(&c.TutorLName)
	fmt.Println("Your ID: ")
	fmt.Scanln(&c.TutorID)
	fmt.Println("--------------------")
	// Call api caller to create a new passenger object
	errMsg := updateClass(c.ClassID, c)
	fmt.Println(errMsg)
	if errMsg != "Success" {
		fmt.Println(errMsg)
	}
}
func deleteClassesMenu() {
	var c Classes
	fmt.Println("-------DELETING CLASS------")
	fmt.Println("Please fill in the following details.")
	fmt.Println("Class Code: ")
	fmt.Scanln(&c.ClassID)
	fmt.Println("--------------------")
	errMsg := deleteClass(c.ClassID)
	if errMsg != "Success" {
		fmt.Println(errMsg)
	}
}
func printClasses() (string, []Classes) {
	//set up url
	url := classURL + "?key" + key
	// Create request
	response, err := http.Get(url)
	//response.Header.Set("Content-Type", "application/json")
	var c []Classes
	var errMsg string
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		// Fetch Request
		data, _ := ioutil.ReadAll(response.Body)
		// Get fail or success msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 422 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
			//since this is a get have to unmarshal (convert json to class)
			json.Unmarshal([]byte(data), &c)
		}
	}
	response.Body.Close()

	return errMsg, c
}
func printClassesMenu() {

	fmt.Println("------ Classes List ------")

	errMsg, c := printClasses()

	if errMsg != "Success" {
		fmt.Println(errMsg)
	} else {
		printSlice(c)
	}
	fmt.Println("--------------------------")
}

func searchClasses(ModuleID string) (string, []Classes) {
	//set up url
	url := "http://localhost:9101/api/v1/class?ModuleID=" + ModuleID + "&key" + key
	// Create request
	response, err := http.Get(url)
	var c []Classes
	var errMsg string
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		// Fetch Request
		data, _ := ioutil.ReadAll(response.Body)
		// Get fail or success msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 422 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
			//since this is a get have to unmarshal (convert json to class)
			json.Unmarshal([]byte(data), &c)
		}
	}
	response.Body.Close()

	return errMsg, c
}
func searchClassesMenu() {
	var c Classes
	fmt.Println("------ Searching for Classes ------")
	fmt.Println("Please fill in the following details.")
	fmt.Println("Module ID: ")
	fmt.Scanln(&c.ModuleID)
	fmt.Println("--------------------------")
	errMsg, cc := searchClasses(c.ModuleID)

	if errMsg != "Success" {
		fmt.Println(errMsg)
	} else {
		printSlice(cc)

	}
	fmt.Println("--------------------------")
}
func main() {
	var option string
	for {
		// Menu
		fmt.Println("------ Welcome to Class Management ------")
		fmt.Println("[1] Create new class")
		fmt.Println("[2] Update classes")
		fmt.Println("[3] Delete classes")
		fmt.Println("[4] Print all classes")
		fmt.Println("[5] Search for classes")
		fmt.Println("[0] Exit application")
		fmt.Println("------------------------------")
		fmt.Println("Enter your option: ")
		fmt.Scanln(&option)

		// Options
		if option == "0" {
			break
		}
		switch option {
		case "1":
			newClassesMenu()
		case "2":
			updateClassesMenu()
		case "3":
			deleteClassesMenu()
		case "4":
			printClassesMenu()
		case "5":
			searchClassesMenu()

		default:
		}
	}

}
