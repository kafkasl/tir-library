package models

import (
	"net/http"
	"sort"
	"testing"
)

func sortReviews(reviews []Review) {
	sort.Slice(reviews, func(i, j int) bool {
		return reviews[i].Text < reviews[j].Text
	})
}

// Inserts reviews into the DB through POST request.
func TestCreateReviews(t *testing.T) {
	expectedReviews := make([]Review, 1)
	isbns := []string{"0-312-86860-X"}

	for i, isbn := range isbns {
		expectedReviews[i] = Review{Text: "Must read for almost everyone", UserId: 1, BookISBN: isbn}
	}

	// Create all expected reviews
	for _, r := range expectedReviews {

		r.Create()
	}
	var gotReviews []Review
	for _, isbn := range isbns {
		bookReviews, statusCode := GetReviews(isbn)
		if statusCode != http.StatusOK {
			t.Errorf("Status code is wrong. Got: %v, Expected: %v\n", statusCode, http.StatusOK)
		}
		gotReviews = append(gotReviews, bookReviews...)
	}

	found := false
	expected := expectedReviews[0]

	for i := 0; i < len(gotReviews); i++ {

		got := gotReviews[i]

		if expected.Text == got.Text && expected.BookISBN == got.BookISBN && expected.UserId == got.UserId {
			found = true
			break
		}

	}
	if ! found {
		t.Errorf("Expected review %v not found. Got: %v", expectedReviews[0], gotReviews)
	}

}

func TestGetReviews(t *testing.T) {

	expectedReviews := make([]Review, 1)
	isbns := []string{"0-312-86860-X"}

	for i, isbn := range isbns {
		expectedReviews[i] = Review{Text: "Must read for almost everyone", UserId: 1, BookISBN: isbn}
	}

	var gotReviews []Review
	for _, isbn := range isbns {
		bookReviews, statusCode := GetReviews(isbn)
		if statusCode != http.StatusOK {
			t.Errorf("Status code is wrong. Got: %v, Expected: %v\n", statusCode, http.StatusOK)
		}
		gotReviews = append(gotReviews, bookReviews...)
	}

	if len(gotReviews) < len(expectedReviews) {
		t.Errorf("Expected at least %v records, got: %v\n%v\n", len(expectedReviews), len(gotReviews), gotReviews)
		return
	}

	sortReviews(expectedReviews)
	sortReviews(gotReviews)

	for i := 0; i < len(expectedReviews); i++ {

		got, expected := gotReviews[i], expectedReviews[i]

		if expected.Text != got.Text || expected.BookISBN != got.BookISBN || expected.UserId != got.UserId {
			t.Errorf("Handler returned wrong review: got \n %v\n %v\n %v\nwant\n %v\n %v\n %v",
				got.Text, got.BookISBN, got.UserId, expected.Text, expected.BookISBN, expected.UserId)
		}

	}

}
