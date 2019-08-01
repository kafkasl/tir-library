package controllers

import (
	"encoding/json"
	"github.com/kafkasl/tir-library/models"
	u "github.com/kafkasl/tir-library/utils"
	"net/http"
)

func CreateAccount(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) // Decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	resp, status := account.Create(false) // create admin without admin rights (can not create admins through API)
	u.Respond(w, resp, status)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	resp, status := models.Login(account.Email, account.Password)
	u.Respond(w, resp, status)
}