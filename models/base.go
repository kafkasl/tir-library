package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"os"
)

var db *gorm.DB

func init() {

	e := godotenv.Load()
	if e != nil {
		fmt.Print("Error loading .env file: ", e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	environment := os.Getenv("environment")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	if environment == "testing" {
		db.Delete(Account{})
		db.Delete(Book{})
		db.Delete(Review{})
		db.Debug().AutoMigrate(&Account{}, &Book{}, &Review{})
	} else {
		db.AutoMigrate(&Account{}, &Book{}, &Review{})
	}

}

func GetDB() *gorm.DB {
	return db
}
