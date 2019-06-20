package controllers

import (
	"github.com/kafkasl/tir-library/models"
	"github.com/kafkasl/tir-library/utils"
	"os"
	"testing"
)

type Body struct {
	Data   []interface{} `json:"data"`
	Message string        `json:"message"`
	Status  bool          `json:"status""`
}

func startup(){

	// Create a user for review inserting
	tester := models.Account{Email: "beta@tester.alph", Password: "omegon"}
	tester.Create(true)

	// Add a book to test for duplicated entries
	book := models.Book{Name: "Ender's Game", Author: "Orson Scott Card", ISBN: "0-312-93208-1"}
	book.Create()

}

func TestMain(m *testing.M) {

	startup()

	// Clean up hook that removes records inserted by the test
	cleaner := utils.DeleteCreatedEntities(models.GetDB())
	defer cleaner()

	code := m.Run()

	os.Exit(code)
}
