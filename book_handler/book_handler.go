package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
	"strings"
	"github.com/gorilla/mux"
	"github.com/yourusername/library-api/models"
)

var mu sync.Mutex

func InitBookRoutes(r *mux.Router) {
	r.HandleFunc("/books", listBooks).Methods("GET")
	r.HandleFunc("/books", addBook).Methods("POST")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books/{id}/checkout", checkoutBook).Methods("PUT")
	r.HandleFunc("/books/{id}/return", returnBook).Methods("PUT")
}

func listBooks(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Books)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if book.Title == "" || book.Author == "" {
		http.Error(w, "Title and author are required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	book.ID = generateID()
	book.CheckedOut = false
	models.Books[book.ID] = book

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	mu.Lock()
	defer mu.Unlock()

	book, exists := models.Books[id]
	if !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func checkoutBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	mu.Lock()
	defer mu.Unlock()

	book, exists := models.Books[id]
	if !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	if book.CheckedOut {
		http.Error(w, "Book already checked out", http.StatusConflict)
		return
	}

	book.CheckedOut = true
	book.DueDate = time.Now().Add(14 * 24 * time.Hour)
	models.Books[id] = book

	models.Loans[id] = models.Loan{
		BookID:   id,
		UserID:   "default-user", // Replace with actual user ID when auth is added
		DueDate:  book.DueDate,
		Returned: false,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func returnBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	mu.Lock()
	defer mu.Unlock()

	book, exists := models.Books[id]
	if !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	if !book.CheckedOut {
		http.Error(w, "Book is not checked out", http.StatusConflict)
		return
	}

	book.CheckedOut = false
	book.DueDate = time.Time{}
	models.Books[id] = book

	loan := models.Loans[id]
	loan.Returned = true
	models.Loans[id] = loan

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func generateID() string {
	return strings.ReplaceAll(time.Now().Format("20060102150405.000000"), ".", "")
}
