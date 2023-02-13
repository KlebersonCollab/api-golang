package main

import (
	"net/http"

	"api-firestore/cmd/routes"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", routes.Home)
	r.HandleFunc("/register", routes.CreateUser)
	r.HandleFunc("/delete{id}", routes.DeleteUser)
	r.HandleFunc("/update{id}", routes.UpdateUser)

	http.ListenAndServe(":8080", r)
}
