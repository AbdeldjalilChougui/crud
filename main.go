package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"log"
	"strconv"
	"crypto/rand"
	"encoding/json"
	"math/big"
)

type Course struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Grade  int     `json:"grade"`
	Author *Author `json:"author"`
}

type Author struct {
	ID     string  `json:"id"`
	Name string `json:"name"`
}


var courses []Course = []Course{
	{ID: "1", Name: "Math", Grade: 5, Author: &Author{ID: "1", Name: "John"}},
	{ID: "2", Name: "Science", Grade: 4, Author: &Author{ID: "2", Name: "Jane"}},
	{ID: "3", Name: "History", Grade: 3, Author: &Author{ID: "3", Name: "Bob"}},
}

func main() {
	fmt.Println("Server Run")

	r := mux.NewRouter()
	r.HandleFunc("/", serveHome)
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/courses/{id}", getOneCourseById).Methods("GET")
	r.HandleFunc("/courses", createCourse).Methods("POST")
	r.HandleFunc("/courses/{id}", updateCourse).Methods("PUT")
	r.HandleFunc("/courses/{id}", deleteCourse).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":4000", r))
}

func (c *Course) GetAuthor() *Author {
	return c.Author
}

func (c *Course) IsEmpty() bool {
	return c.Name == ""
}

func serveHome (w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1> Welcome to Home </h1>"))
}

func getAllCourses (w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all courses")
	w.Header().Set("Content-type", "application/json")

	fmt.Println(courses[0].ID)
	json.NewEncoder(w).Encode(courses)
}

func getOneCourseById (w http.ResponseWriter, r *http.Request) {	
	fmt.Println("Get one course by id")
	w.Header().Set("Content-type", "application/json")

	params := mux.Vars(r)

	for _, course := range courses {
		if course.ID == params["id"] {
			json.NewEncoder(w).Encode(course)
			return
		}
	}

	json.NewEncoder(w).Encode("No Course found with this id")
	return
}

func createCourse (w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create course")
	w.Header().Set("Content-type", "application/json")

	if r.Body == nil {
		json.NewEncoder(w).Encode("No body sent in request")
	}

	var course Course
	json.NewDecoder(r.Body).Decode(&course)

	course.Name = "name"
	if course.IsEmpty() {
		json.NewEncoder(w).Encode("Course name is required")
		return
	}
	
	id, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		json.NewEncoder(w).Encode("Failed to generate random ID")
		return
	}
	course.ID = strconv.FormatInt(id.Int64(), 10)
	
	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)
	return
}

func updateCourse (w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update course")
	w.Header().Set("Content-type", "application/json")

	params := mux.Vars(r)

	for index, course := range courses {
		if course.ID == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			var course Course
			_ = json.NewDecoder(r.Body).Decode(&course)
			
			course.ID = params["id"] 
			course.Name = params["name"] 
			grade, err := strconv.Atoi(params["grade"])
			if err != nil {
				json.NewEncoder(w).Encode("Invalid grade value")
				return
			}
			course.Grade = grade
			courses = append(courses, course)

			json.NewEncoder(w).Encode(course)
			return
		}
	}
}


func deleteCourse (w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete course")
	w.Header().Set("Content-type", "application/json")

	params := mux.Vars(r)

	for index, course := range courses {
		if course.ID == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(courses)
	return 
}
