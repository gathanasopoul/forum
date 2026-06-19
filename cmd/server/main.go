package main

import (
	"fmt"
	"net/http"

	"forum/internal/database"
	"forum/internal/handlers"
)

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
	http.HandleFunc("/", h.HomePage)
	http.HandleFunc("/register", h.RegisterPage)
	http.HandleFunc("/login", h.LoginPage)
	http.HandleFunc("/logout", h.Logout)
	http.HandleFunc("/create-post", h.CreatePostPage)
	http.HandleFunc("/post", h.ViewPostPage)
	http.HandleFunc("/comment/create", h.CreateComment)
	http.HandleFunc("/like", h.LikePost)

	fmt.Println("Server running on :8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
