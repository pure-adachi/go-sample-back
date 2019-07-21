package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type Model struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index" json:"-"`
}

type Todo struct {
	Model
	Name  string
	State int
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db := GetDBConn()

	todos := []Todo{}

	db.Find(&todos)

	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	todo := Todo{}
	_ = json.NewDecoder(r.Body).Decode(&todo)
	todo.State = 1

	db := GetDBConn()
	db.Create(&todo)

	json.NewEncoder(w).Encode(todo)
}

// TODO
func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	params := mux.Vars(r)
	todo := Todo{}

	db := GetDBConn()
	db.First(&todo, params["id"])
	_ = json.NewDecoder(r.Body).Decode(&todo)

	db.Save(&todo)

	json.NewEncoder(w).Encode(todo)
}

// TODO
func deleteodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	todo := Todo{}

	db := GetDBConn()
	db.Delete(&todo, params["id"])

	json.NewEncoder(w).Encode(Todo{})
}

func GetDBConn() *gorm.DB {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=go-test sslmode=disable")
	if err != nil {
		panic(err)
	}

	db.LogMode(true)
	return db
}

func main() {
	db := GetDBConn()
	db.AutoMigrate(&Todo{})

	r := mux.NewRouter()

	r.HandleFunc("/api/todos", getTodos).Methods("GET")
	r.HandleFunc("/api/todos", createTodo).Methods("POST")
	r.HandleFunc("/api/todos/{id}", updateTodo).Methods("PUT")
	r.HandleFunc("/api/todos/{id}", deleteodo).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
