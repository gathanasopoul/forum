package main

import (
	"fmt"
	"net/http"

	"forum/internal/database"
	"forum/internal/handlers"
)

// Home route
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Forum Home")
}

func main() {

	// Connect database
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Database connected")

	// Run migrations
	err = database.RunMigrations(db)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database initialized")

	// Routes
	h := handlers.New(db)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/register", h.RegisterPage)
	http.HandleFunc("/login", h.LoginPage)
	http.HandleFunc("/logout", h.Logout)

	fmt.Println("Server running on :8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
