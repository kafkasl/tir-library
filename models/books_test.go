package models

import (
	"net/http"
	"sort"
	"testing"
)

func sortBooks(books []Book) {
	sort.Slice(books, func(i, j int) bool {
		return books[i].Name < books[j].Name
	})
}

// Inserts books into the DB through POST request.
func TestCreateBooks(t *testing.T) {

	expectedBooks := make([]Book, 1)
	expectedBooks[0] = Book{Name: "Shadow of the Hegemon", Author: "Orson Scott Card", ISBN: "0-312-87651-3"}

	// Insert all expected books through POST requests
	for _, b := range expectedBooks {
		b.Create()
	}

	gotBooks, statusCode := GetBooks()

	if statusCode != http.StatusOK {
		t.Errorf("Status code is wrong. Got: %v, Expected: %v\n", statusCode, http.StatusOK)
	}

	sortBooks(expectedBooks)
	sortBooks(gotBooks)

	found := false
	for i := 0; i < len(gotBooks); i++ {

		expected := expectedBooks[0]
		got := gotBooks[i]

		if expected.Name == got.Name && expected.Author == got.Author && expected.ISBN == got.ISBN {
			found = true
			break
		}

	}
	if ! found {
		t.Errorf("Expected book %v not found. Got books: %v", expectedBooks[0], gotBooks)
	}

}

func TestGetBooks(t *testing.T) {

	expectedBooks := make([]Book, 2)
	expectedBooks[0] = Book{Name: "Ender's Game", Author: "Orson Scott Card", ISBN: "0-312-93208-1"}
	expectedBooks[1] = Book{Name: "Ender's Shadow", Author: "Orson Scott Card", ISBN: "0-312-86860-X"}

	gotBooks, statusCode := GetBooks()

	if statusCode != http.StatusOK {
		t.Errorf("Status code is wrong. Got: %v, Expected: %v\n", statusCode, http.StatusOK)
	}

	if len(gotBooks) < len(expectedBooks) {
		t.Errorf("Expected at least %v records, got: %v\n%v\n", len(expectedBooks), len(gotBooks), gotBooks)
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
