package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/kafkasl/tir-library/models"
	u "github.com/kafkasl/tir-library/utils"
	"net/http"
)

var CreateBook = func(w http.ResponseWriter, r *http.Request) {

	book := &models.Book{}
	err := json.NewDecoder(r.Body).Decode(book)
	if err != nil {
		message := fmt.Sprintf("Error while decoding request body: %v\n", err)
		u.Respond(w, u.Message(false, message), http.StatusBadRequest)
		return
	}

	resp, status := book.Create()
	u.Respond(w, resp, status)
}

var GetBooks = func(w http.ResponseWriter, r *http.Request) {

	data, status := models.GetBooks()
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp, status)
}
