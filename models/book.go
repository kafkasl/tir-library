package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	u "github.com/kafkasl/tir-library/utils"
	"net/http"
)

type Book struct {
	gorm.Model
	Name   string `json:"name"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
}

func unique(isbn string) (bool, interface{}) {
	temp := &Book{}

	err := GetDB().Table("books").Where("isbn = ?", isbn).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if temp.ISBN != "" {
		return false, nil
	}
	return true, nil
}

func valid(isbn string) bool {
	// this function should check that the ISBN is correctly formatted.
	return true
}

/*
 This struct function validate the required parameters sent through the http request body

returns message and true if the requirement is met
*/
func (book *Book) Validate() (map[string]interface{}, bool) {

	if book.Name == "" {
		return u.Message(false, "Book name not present."), false
	}

	if book.Author == "" {
		return u.Message(false, "Book author not present."), false
	}

	unique, err := unique(book.ISBN)
	if err != nil {
		return u.Message(false, fmt.Sprintf("Connection error checking for duplicated ISBN: %v. Please retry", err)), false
	}
	if ! unique {
		return u.Message(false, "Book ISBN is already registered"), false
	}
	if ! valid(book.ISBN) {
		return u.Message(false, "Book ISBN format is not valid."), false
	}

	fmt.Println("Succesful book validation")

	// All the required parameters are present
	return u.Message(true, "Book has been created"), true
}

func (book *Book) Create() (map[string]interface{}, int) {

	resp, ok := book.Validate()

	if !ok {
		return resp, http.StatusBadRequest
	}

	GetDB().Create(book)

	resp["book"] = book
	return resp, http.StatusCreated
}

func GetBook(isbn uint) (*Book) {

	book := &Book{}
	err := GetDB().Table("books").Where("isbn = ?", isbn).First(book).Error
	if err != nil {
		return nil
	}
	return book
}

func GetBooks() ([]*Book, int) {

	books := make([]*Book, 0)
	db := GetDB()
	if db == nil {
		fmt.Println("Nil DB")
		return nil, http.StatusInternalServerError
	}

	table := db.Table("books")
	err := table.Find(&books).Error

	if err != nil {
		fmt.Println(err)
		return nil, http.StatusInternalServerError
	}

	return books, http.StatusOK
}
