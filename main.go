package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kafkasl/tir-library/app"
	"github.com/kafkasl/tir-library/controllers"
	"net/http"
	"os"
)

func main() {

	router := mux.NewRouter()

	// Accounts
	router.HandleFunc("/api/user/signup", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")

	// Books
	router.HandleFunc("/api/books", controllers.CreateBook).Methods("POST")
	router.HandleFunc("/api/books", controllers.GetBooks).Methods("GET")

	// Book Reviews
	router.HandleFunc("/api/books/{book_isbn}/reviews", controllers.GetReviewsFor).Methods("GET")
	router.HandleFunc("/api/books/{book_isbn}/reviews", controllers.CreateReview).Methods("POST")

	router.Use(app.JwtAuthentication) //  attach JWT auth middleware


	// Port set by environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // localhost dev
	}

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
