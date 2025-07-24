package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Book represents the model
type Book struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&Book{})

	r := mux.NewRouter()
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")

	log.Println("API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book
	db.Find(&books)
	json.NewEncoder(w).Encode(books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var input Book
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	// Validation: title and author must be provided
	if input.Title == "" || input.Author == "" {
		http.Error(w, "Missing required fields: title and author", http.StatusBadRequest)
		return
	}

	result := db.Create(&input)
	if result.Error != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}