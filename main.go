package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	// "path/filepath"
	"runtime"
	"time"
)

// global variables

// records variable
var records = []Record{}

// Record struct

type Record struct {
	UserInstance
	Text   string
	Id     int
	Status string
	Saved  bool
}

// User Interface

type User interface {
	welcomeUser()
	loadUserRecords(*Record)
	addRecord(string, *Record) Record
	updateRecord(int, string) Record
	deleteRecord(int)
	viewAllRecords(id int) []Record
}

// User struct

type UserInstance struct {
	Name   string
	UserId int
}

// main functions

// Introduction message
// Welcome and get the user's information
func (user *UserInstance) welcomeUser() {
	fmt.Printf("Welcome to the program %s ,", user.Name)
	fmt.Println("Please select an option")
	fmt.Println("1. Add a new record")
	fmt.Println("2. Update a record")
	fmt.Println("3. Delete a record")
	fmt.Println("4. View all records")
	fmt.Println("5. Exit")

}

// Add a new record
func (u UserInstance) addRecord(collection *[]Record, record string) Record {

	if record == "" {
		fmt.Println("Entry cannot be empty")
		fmt.Scan(&record)
	}
	defer fmt.Println("Record added successfully")
	fmt.Println("Adding a new record")
	fmt.Println("Enter the record text")
	newRecord := Record{UserInstance: UserInstance{Name: u.Name, UserId: u.UserId}, Text: record, Id: randomIdGenerator(), Status: "incomplete", Saved: false}
	*collection = append(records, newRecord)

	return newRecord
}

// load user's records data tp app
func (u UserInstance) loadUserRecords(collection *[]Record) {
	data, _ := loadRecords()

	// var userRecords = []Record{}

	// var recordsPoint = &records

	for _, record := range data {
		if record.UserInstance.UserId == u.UserId {

			*collection = append(records, record)
		}
	}

	// *collection = append(records, userRecords...)

	fmt.Println("from load user recors", records)

}

// Update a record
func (u UserInstance) updateRecord(id int, record string) Record {
	defer fmt.Println("Record updated successfully")
	fmt.Println("Updating a record")
	fmt.Println("Enter the record text")
	records[id].Text = record

	return records[id]
}

// Delete a record
func (u UserInstance) deleteRecord(collection *[]Record, id int) {
	defer fmt.Println("Record deleted successfully")
	fmt.Println("Deleting a record")
	oldRecords, err := loadRecords()

	if err != nil {

		oldRecords = deleteRecord(oldRecords, id, u.UserId)
		err := os.Remove("./records.json")

		if err != nil {
			fmt.Println("Error deleting records file:", err)
			return
		}

		saveRecords(oldRecords)
	}

	*collection = deleteRecord(records, id, u.UserId)

}

// View all records
func (u UserInstance) viewAllRecords() []Record {

	fmt.Println("Displaying all your records")

	return records
}

func (u UserInstance) markTaskComplete(id int) {
	fmt.Println("Marking task as complete", id)
}

// find and delete entry from a slice

func deleteRecord(collection []Record, id int, userId int) []Record {
	for index, value := range collection {
		if value.Id == id && value.UserId == userId {
			collection = append(collection[:index], collection[index+1:]...)
		}
	}
	return collection
}

func randomIdGenerator() int {
	return rand.Intn(100)
}

// mark records as saved before writing to file
func markRecordsAsSaved(collection []Record) []Record {
	updatedCollection := []Record{}

	if len(collection) > 0 {
		for _, value := range collection {

			value.Saved = true
			updatedCollection = append(updatedCollection, value)
		}
	}

	return updatedCollection
}

// check for repeated records

func filterRepeatedRecords(collection1 []Record, collection2 []Record) []Record {
	seen := make(map[string]bool)
	saved := make(map[string]bool)

	for _, record1 := range collection1 {
		key := fmt.Sprintf("%v-%v", record1.UserId, record1.Id)
		seen[key] = true
		saved[key] = record1.Saved
	}

	// var filteredCollection []Record
	for index, record2 := range collection2 {
		key := fmt.Sprintf("%v-%v", record2.UserId, record2.Id)
		if !seen[key] && saved[key] {
			collection2 = append(collection2[:index], collection2[index+1:]...)
		}
	}
	return collection2
}

// persisting the records data

// save data to file system
func saveRecords(collection []Record) {
	if len(collection) == 0 {
		fmt.Println("No records to save.")
		return
	}

	fmt.Print("Saving data to file system... ")
	defer fmt.Println("Done.")

	// get old records
	oldRecords, err := loadRecords()

	if err != nil {
		fmt.Println("Error loading old records or old records not found")
	}

	// //merge old records with new records
	for _, record := range collection {
		if !record.Saved {
			oldRecords = append(oldRecords, record)
		}
	}

	// check for reapeated fields
	oldRecords = filterRepeatedRecords(collection, oldRecords)

	fmt.Println(oldRecords)

	//markRecords as saved
	oldRecords = markRecordsAsSaved(oldRecords)

	// Marshal records into JSON
	recordsJson, err := json.MarshalIndent(oldRecords, "", "  ")
	if err != nil {
		fmt.Println("\nError marshalling records:", err)
		return
	}

	_, err = loadRecords()

	if err != nil {
		file, err := os.OpenFile("./records.json", os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			fmt.Println("Error opening or creating records.json:", err)
			return
		}
		defer file.Close()
	}

	// Write to file
	err = os.WriteFile("records.json", recordsJson, 0644)
	if err != nil {
		fmt.Println("\nError writing to file:", err)
		return
	}
}

// loading data from file system

func loadRecords() ([]Record, error) {
	data, err := os.ReadFile("records.json")
	if err != nil {
		return nil, fmt.Errorf("error loading records: %w", err)
	}

	var allRecords []Record
	err = json.Unmarshal(data, &allRecords)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling records: %w", err)
	}

	fmt.Println("Records loaded successfully from file system")
	return allRecords, nil
}

// check if the user exists

type Info struct {
	exists bool
	UserId int
}

func checkExistingUser(name string) Info {
	allRecords, err := loadRecords()

	if err != nil {
		fmt.Println("Error loading records")
		return Info{exists: false, UserId: 0}
	}

	for _, value := range allRecords {
		if value.UserInstance.Name == name {
			return Info{exists: true, UserId: value.UserInstance.UserId}
		}
	}
	return Info{exists: false, UserId: 0}
}

// app runner function

func runApp() {
	//set a seed for the random number generator
	rand.NewSource(time.Now().UnixNano())

	// create a new user
	var user UserInstance
	var name string
	fmt.Println("Enter your user name")
	fmt.Scanln(&name)

	// check if the user exists

	if checkExistingUser(name).exists {
		userInfo := checkExistingUser(name)
		fmt.Println("Welcome back", name)
		user = UserInstance{Name: name, UserId: userInfo.UserId}
		user.loadUserRecords(&records)
	} else {
		fmt.Println("Welcome", name)
		user = UserInstance{Name: name, UserId: rand.Intn(100)}
	}

	// get the operating system
	osRuntime := runtime.GOOS

	// set the operating system
	os.Setenv("GOOS", osRuntime)

	// print the operating system
	fmt.Printf("You are running on a %s machine", runtime.GOOS)

	var choice int

	for choice != 5 {
		// welcome the user
		user.welcomeUser()

		// get the user's choice
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			var record string
			fmt.Println("Enter the record text")
			fmt.Scanln(&record)
			user.addRecord(&records, record)
		case 2:
			var id int
			fmt.Println("Enter the record id")
			fmt.Scanln(&id)
			var record string
			fmt.Println("Enter the record text")
			fmt.Scanln(&record)
			user.updateRecord(id, record)
		case 3:
			var id int
			fmt.Println("Enter the record id")
			fmt.Scanln(&id)
			user.deleteRecord(&records, id)
		case 4:
			allRecords := user.viewAllRecords()
			fmt.Println(allRecords)
		case 5:
			saveRecords(records)
			os.Exit(0)
		default:
			fmt.Println("Invalid choice")

		}
	}

}
func main() {

	runApp()

}
