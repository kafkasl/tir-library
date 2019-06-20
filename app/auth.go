package app

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kafkasl/tir-library/models"
	u "github.com/kafkasl/tir-library/utils"
	"net/http"
	"os"
	"strings"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		requestMethod := r.Method

		notAuth := []string{"/api/user/signup",
			"/api/user/login"}

		// Signup and Login do not require authentication
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Similarly, book GET methods do not require auth either
		if strings.HasPrefix(requestPath, "/api/books") && requestMethod == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") // Grab the token from the header

		if tokenHeader == "" { //  Token is missing, returns with error code 403 Unauthorized
			response = u.Message(false, "Missing auth token")
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response, http.StatusForbidden)
			return
		}

		splitted := strings.Split(tokenHeader, " ") // JWT usual format `Bearer {token-body}`
		// check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid/Malformed auth token")
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response, http.StatusForbidden)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			response = u.Message(false, "Malformed authentication token")
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response, http.StatusForbidden)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			response = u.Message(false, "Token is not valid.")
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response, http.StatusForbidden)
			return
		}

		//  Auth. successful, proceed with the request and set the caller to the user retrieved from the parsed token
		fmt.Sprintf("User %v", tk.UserId) //Useful for monitoring
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	});
}
