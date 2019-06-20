package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	u "github.com/kafkasl/tir-library/utils"
	"net/http"
)

type Review struct {
	gorm.Model
	Text     string `json:"text"`
	UserId   uint   `json:"user_id"`   //The user that this review belongs to
	BookISBN string `json:"book_isbn"` //The book that this review belongs to
}

func existsISBN(isbn string) (bool, interface{}) {
	temp := &Book{}

	err := GetDB().Table("books").Where("isbn = ?", isbn).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if temp.ISBN != "" {
		return true, nil
	}
	return false, nil
}

// Validate checks that the struct function pass contains the required parameters sent through the http request body.
// Returns message and true if the requirement is met
func (review *Review) Validate() (map[string]interface{}, bool) {
	if review.Text == "" {
		return u.Message(false, "Empty reviews are not allowed here."), false
	}

	exists, err := existsISBN(review.BookISBN)
	if err != nil {
		return u.Message(false, fmt.Sprintf("Connection error checking for existing ISBN: %v. Please retry", err)), false
	}

	if ! exists {
		return u.Message(false, fmt.Sprintf("Book with ISBN %v not found.", review.BookISBN)), false
	}

	if review.UserId <= 0 {
		return u.Message(false, "User is not recognized"), false
	}

	// All the required parameters are present
	return u.Message(true, "success"), true
}

func (review *Review) Create() (map[string]interface{}, int) {

	if resp, ok := review.Validate(); !ok {
		return resp, http.StatusBadRequest
	}

	GetDB().Create(review)

	resp := u.Message(true, "Review has been created")
	resp["review"] = review
	return resp, http.StatusCreated
}

func GetReviews(isbn string) ([]*Review, int) {

	reviews := make([]*Review, 0)
	err := GetDB().Table("reviews").Where("book_isbn = ?", isbn).Find(&reviews).Error
	if err != nil {
		fmt.Println(err)
		return nil, http.StatusInternalServerError
	}

	return reviews, http.StatusOK
}
