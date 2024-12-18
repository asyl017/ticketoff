package main

import (
	"fmt"
	"net/http"
)

func home_page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the home page")
}
func contacts_page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Contacts page")
}
func handleRequest() {
	http.HandleFunc("/", home_page)
	http.HandleFunc("/contacts/", contacts_page)
	http.ListenAndServe(":8080", nil)
}
func main() {
	handleRequest()
}
