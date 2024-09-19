package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

type User struct {
	Name string `json:"name"`
}

var userCache = make(map[int]User)

var cacheMutex sync.RWMutex

func main() {

	db, err := sql.Open("sqlite3", "C:/sqlite/database.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)

	mux.HandleFunc("POST /users", createUser)
	mux.HandleFunc("GET /users/{id}", getUser)
	mux.HandleFunc("DELETE /users/{id}", deleteUser)
	fmt.Println("Server listening to : 8080")
	http.ListenAndServe(":8080", mux)

	defer db.Close()

}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if _, ok := userCache[id]; !ok {
		http.Error(w, "not Found", http.StatusBadRequest)
	}
	cacheMutex.Lock()
	delete(userCache, id)
	cacheMutex.Unlock()
	w.WriteHeader(http.StatusNoContent)
}
func getUser(w http.ResponseWriter, r *http.Request) {
	// Retrieve the ID from the request URL path.
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}

	// Open a connection to the SQLite database.
	db, err := sql.Open("sqlite3", "C:/sqlite/database.db") // Replace with your actual database path.
	if err != nil {
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query the database for the user with the given ID.
	row := db.QueryRow("SELECT id FROM users WHERE id = ?", id)

	var user User
	err = row.Scan(&id)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}

	// Marshal the user data to JSON.
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	// Write the response.
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User

	// Decode the JSON request body into the User struct.
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate that the Name field is not empty.
	if user.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Open a connection to the SQLite database.
	db, err := sql.Open("sqlite3", "C:/sqlite/database.db")
	if err != nil {
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Insert the new user into the database.
	_, err = db.Exec("INSERT INTO users (name) VALUES (?)", user.Name)
	if err != nil {
		http.Error(w, "Error inserting user into database", http.StatusInternalServerError)
		return
	}

	// Respond with a status indicating that the user was created successfully.
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User created successfully")

}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")

}
