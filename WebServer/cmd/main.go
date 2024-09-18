package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type User struct {
	Name string `json:"name"`
}

var userCache = make(map[int]User)

var cacheMutex sync.RWMutex

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)

	mux.HandleFunc("POST /users", createUser)
	mux.HandleFunc("GET /users/{id}", getUser)
	fmt.Println("Server listening to : 8080")
	http.ListenAndServe(":8080", mux)

}
func getUser(w http.ResponseWriter, r *http.Request) {

}
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return

	}
	cacheMutex.Lock()
	userCache[len(userCache)+1] = user
	cacheMutex.Unlock()
	w.WriteHeader(http.StatusNoContent)

}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")

}
