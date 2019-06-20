package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kafkasl/tir-library/models"
	u "github.com/kafkasl/tir-library/utils"
	"net/http"
)

var CreateReview = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint) // Grab the id of the user that send the request
	book_isbn := mux.Vars(r)["book_isbn"]

	review := &models.Review{}

	err := json.NewDecoder(r.Body).Decode(review)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"), http.StatusBadRequest)
		return
	}

	review.UserId = user
	review.BookISBN = book_isbn

	resp, status := review.Create()
	u.Respond(w, resp, status)
}

var GetReviewsFor = func(w http.ResponseWriter, r *http.Request) {

	book_isbn := mux.Vars(r)["book_isbn"]
	data, status := models.GetReviews(book_isbn)

	resp := u.Message(true, "success")
	resp["data"] = data

	u.Respond(w, resp, status)
}
