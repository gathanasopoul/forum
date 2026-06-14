package main

import (
	"fmt"
	"net/http"
)

// Home route
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Forum Home")
}

func main() {

	http.HandleFunc("/", homeHandler)

	fmt.Println("Server running on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
