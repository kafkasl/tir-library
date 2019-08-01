package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/kafkasl/tir-library/models"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
)

func sortAccounts(Accounts []models.Account) {
	sort.Slice(Accounts, func(i, j int) bool {
		return Accounts[i].ID < Accounts[j].ID
	})
}

// Inserts Accounts into the DB through POST request.
func insertAccounts(t *testing.T, users []models.Account) {

	// Insert all expected Accounts through POST requests
	for _, u := range users {

		locJson, err := json.Marshal(u)

		req, err := http.NewRequest("POST", "/api/users/new", bytes.NewBuffer(locJson))
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateAccount)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	}

}

// getAccounts test that expectedAccounts are present in DB.
func getAccounts(t *testing.T, expectedAccounts []models.Account) {

	gotAccounts := models.GetUsers()

	if len(gotAccounts) < len(expectedAccounts) {
		t.Errorf("Expected at least %v records, got: %v\n%v\n", len(expectedAccounts), len(gotAccounts), gotAccounts)
		return
	}

	sortAccounts(expectedAccounts)
	sortAccounts(gotAccounts)

	for i := 0; i < len(expectedAccounts); i++ {

		got, expected := gotAccounts[i], expectedAccounts[i]
		
		if expected.Admin != got.Admin || expected.Email != got.Email || expected.Token != got.Token {
			t.Errorf("Handler returned wrong Account: got \n %v\n %v\n %v\n "+
				"want\n %v\n %v\n %v\n",
				got.Email, got.Token, got.Admin, expected.Email, expected.Token, expected.Admin)
		}

	}

}

func TestAccounts(t *testing.T) {

	expected := make([]models.Account, 2)
	expected[0] = models.Account{Email: "arthur@heart.gold", Password: "424242"}
	expected[1] = models.Account{Email: "trurl@cyberiad.cosmos", Password: "Klapaucius"}

	insertAccounts(t, expected)

	getAccounts(t, expected)

}
