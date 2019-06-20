package models

import (
	"testing"
)

func TestBook(t *testing.T) {

	book := Book{Name: "Ender's Game", Author: "Orson Scott Card", ISBN: "0-312-93208-1"}

	if r, valid := book.Validate(); ! valid {
		t.Errorf("Book information is not valid and should be. Reason: %v\nBook: %v\n", r["message"], book)
	}

	if r := book.Create(); ! r["status"].(bool) {
		t.Errorf("Error creating book. Reason: %v\nBook: %v\n", r["message"], book)

	}

	if r, valid := book.Validate(); valid {
		t.Errorf("Book information is valid and should be because the same book was already inserted. Reason: %v\nBook: %v\n", r["message"], book)
	}

	if r := book.Create(); r["status"].(bool) {
		t.Errorf("Not error while creating duplicated book. Reason: %v\nBook: %v\n", r["message"], book)

	}

}
