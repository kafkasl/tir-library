package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kafkasl/tir-library/models"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
)

func sortReviews(reviews []models.Review) {
	sort.Slice(reviews, func(i, j int) bool {
		return reviews[i].Text < reviews[j].Text
	})
}

// Inserts reviews into the DB through POST request.
func insertReviews(t *testing.T, reviews []models.Review) {

	// Insert all expected reviews through POST requests
	for _, b := range reviews {

		locJson, err := json.Marshal(b)

		req, err := http.NewRequest("POST", "/api/reviews", bytes.NewBuffer(locJson))
		q := req.URL.Query()
		q.Add("user", "1") // add mockup user ID
		req.URL.RawQuery = q.Encode() // Encode and assign back to the original query.

		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateReview)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	}

}

// Retrieve all the DB reviews using GET and test that they are equal to expectedReviews.
func getReviews(t *testing.T, expectedReviews []models.Review) {

	var gotReviews []models.Review
	//make([]models.Review)

	books, _ := models.GetBooks()
	for _, b := range books {
		url := fmt.Sprintf("/api/books/%v/reviews", b.ISBN)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetReviewsFor)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		data := Body{}

		err = json.NewDecoder(rr.Body).Decode(&data)
		if err != nil {
			t.Errorf("Decoding err: %v", err)
		}

		for _, r := range data.Data {
			gotReviews = append(gotReviews, r.(models.Review))
		}
	}
	if len(gotReviews) < len(expectedReviews) {
		t.Errorf("Expected %v records, got: %v\n", len(expectedReviews), len(gotReviews))
		return
	}

	sortReviews(expectedReviews)
	sortReviews(gotReviews)

	for i := 0; i < len(expectedReviews); i++ {

		got, expected := gotReviews[i], expectedReviews[i]

		if expected.Text != got.Text || expected.BookISBN != got.BookISBN || expected.UserId != got.UserId {
			t.Errorf("Handler returned wrong review: got \n %v\n %v\n %v\nwant\n %v\n %v\n %v",
				got.UserId, got.BookISBN, got.Text, expected.UserId, expected.BookISBN, expected.Text)
		}

	}

}

func TestReviews(t *testing.T) {

	books, _ := models.GetBooks()

	expected := make([]models.Review, 2*len(books))

	users := models.GetUsers()

	if len(users) < 1 {
		t.Errorf("No users on DB to be used as reviewer.")
	}

	reviewer := users[0]


	for idx, b := range books {
		expected[idx] = models.Review{Text: fmt.Sprintf("R1: %v is awesome.", b.Name), BookISBN: b.Name, UserId: reviewer.ID}
		expected[idx+1] = models.Review{Text: fmt.Sprintf("R2: %v is awesome.", b.Name), BookISBN: b.Name,  UserId: reviewer.ID}
	}

	insertReviews(t, expected)

	getReviews(t, expected)

}
