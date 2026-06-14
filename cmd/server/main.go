package main

import (
	"fmt"
	"net/http"

	"forum/internal/database"
)

// Home route
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Forum Home")
}

func main() {

	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Database connected")

	err = database.RunMigrations(db)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database initialized")

	http.HandleFunc("/", homeHandler)

	fmt.Println("Server running on :8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
