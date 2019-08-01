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


func sortBooks(books []models.Book) {
	sort.Slice(books, func(i, j int) bool {
		return books[i].Name < books[j].Name
	})
}

// Inserts books into the DB through POST request.
func insertBooks(t *testing.T, books []models.Book) {

	// Insert all expected books through POST requests
	for _, b := range books {
		fmt.Printf("Inserting %v\n", b)

		locJson, err := json.Marshal(b)

		req, err := http.NewRequest("POST", "/api/books", bytes.NewBuffer(locJson))
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateBook)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	}


}

// Retrieve all the DB books using GET and test that they are equal to expectedBooks.
func getBooks(t *testing.T, expectedBooks []models.Book) {
	req, err := http.NewRequest("GET", "/api/books", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetBooks)
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


	var gotBooks []models.Book
	for _, b := range data.Data {
		book, _ := b.(models.Book)
		gotBooks = append(gotBooks, book)
	}

	if len(gotBooks) < len(expectedBooks) {
		t.Errorf("Expected at least %v records, got: %v\n", len(expectedBooks), len(gotBooks))
		return
	}

	sortBooks(expectedBooks)
	sortBooks(gotBooks)

	for i := 0; i < len(expectedBooks); i++ {

		got, expected := gotBooks[i], expectedBooks[i]

		if expected.Name != got.Name || expected.Author != got.Author || expected.ISBN != got.ISBN {
			t.Errorf("Handler returned wrong book: got \n %v\n %v\n %v\nwant\n %v\n %v\n %v",
				got.Name, got.Author, got.ISBN, expected.Name, expected.Author, expected.ISBN)
		}

	}

}

func TestBooks(t *testing.T) {

	expected := make([]models.Book, 1)
	expected[0] = models.Book{Name: "Ender's Game", Author: "Orson Scott Card", ISBN: "0-312-93208-1"}
	insertBooks(t, expected)

	expected[1] = models.Book{Name: "Ender's Shadow", Author: "Orson Scott Card", ISBN: "0-312-86860-X"}
	insertBooks(t, expected)

	getBooks(t, expected)

	fmt.Println(models.GetBooks())

}