package models

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	u "github.com/kafkasl/tir-library/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strings"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//a struct to rep user account
type Account struct {
	gorm.Model
	Email string `json:"email"`
	Password string `json:"password"`
	Admin bool `json:"admin"`
	Token string `json:"token";sql:"-"`
}

// Validate incoming user details.
func (account *Account) Validate() (map[string] interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password should have at least 6 chars."), false
	}

	//Email must be unique
	temp := &Account{}

	// Check for errors and duplicate emails
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user"), false
	}

	return u.Message(false, "Requirement passed"), true
}

// Create generates a password, and JWToken for account and stores it in the DB.
func (account *Account) Create(admin bool) (map[string] interface{}, int) {

	if resp, ok := account.Validate(); !ok {
		return resp, http.StatusBadRequest
	}

	account.Admin = admin

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error"), http.StatusInternalServerError
	}

	// Create new JWT token for the newly registered account
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" // delete password before returning info

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response, http.StatusCreated
}

// Login takes an email and password and, if it matches a registered user, returns its information with the JWT.
// Status return parameter indicates if the log in was successful.
func Login(email, password string) (map[string]interface{}, int) {

	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found"), http.StatusForbidden
		}
		return u.Message(false, "Connection error. Please retry"), http.StatusInternalServerError
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again"), http.StatusForbidden
	}
	// Logged In
	account.Password = ""

	// Create JWT token
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString // Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp, http.StatusOK
}

// GetUser returns the information for user with id u or nil if the user does not exists.
// User password field, when user is found, is blank.
func GetUser(u uint) *Account {

	acc := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}

// GetUsers returns all the registered users.
func GetUsers() ([]Account) {

	users := make([]Account, 0)
	db := GetDB()
	if db == nil {
		fmt.Println("Nil DB")
		return nil
	}

	table := db.Table("accounts")
	err := table.Find(&users).Error

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return users
}
