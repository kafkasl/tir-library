package models

import (
	"database/sql"
	"fmt"
	"gopkg.in/testfixtures.v2"
	"log"
	"os"
	"testing"
)

var (
	sqlDB    *sql.DB
	fixtures *testfixtures.Context
)

type Body struct {
	Data    []interface{} `json:"data"`
	Message string        `json:"message"`
	Status  bool          `json:"status""`
}

func TestMain(m *testing.M) {

	var err error

	sqlDB = GetDB().DB()
	if db == nil {
		fmt.Println("Nil DB")
		return
	}
	// creating the context that hold the fixtures
	// see about all compatible databases in this page below
	fixtures, err = testfixtures.NewFolder(sqlDB, &testfixtures.PostgreSQL{}, "../fixtures")
	if err != nil {
		log.Fatal(err)
	}

	prepareTestDatabase()

	os.Exit(m.Run())

}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}
