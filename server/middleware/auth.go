package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// We need a custom local string type to avoid variable naming collisions inside the request context
type contextKey string

// Use the contextKey type to define the UserIDKey
const UserIDKey contextKey = "userID"

// Helper function that extracts the userID from the request's context, will be used by other handlers to verify authenticity
func GetUserID(request *http.Request) (int, bool) {
	userID, isOK := request.Context().Value(UserIDKey).(int)

	return userID, isOK
}

// The Middleware wrapper that checks the Authentication header, verifies the validity of the token
// If valid, extracts the userID and adds it to the request's context
// If not, it rejects the request for invalid tokens
func RequireAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		// Extract the value of the Auth header
		authValue := request.Header.Get("Authorization")
		if authValue == "" {
			http.Error(response, "Authorization Header Missing", http.StatusUnauthorized)
			return
		}

		// Split it to extract the JWT Token
		authParts := strings.Split(authValue, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			http.Error(response, "Invalid Authorization Header", http.StatusUnauthorized)
			return
		}

		stringToken := authParts[1]
		// Validate the algorithm used to sign the JWT Token
		parsedToken, err := jwt.Parse(stringToken, func(tkn *jwt.Token) (any, error) {
			_, ok := tkn.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !parsedToken.Valid {
			http.Error(response, "Invalid Or Expired JWT Token", http.StatusUnauthorized)
			return
		}

		// Extract the claims map
		decodedClaims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(response, "Invalid Claims", http.StatusUnauthorized)
			return
		}

		// Extract the float64 version of the id from the claims map
		floatUserID, ok := decodedClaims["id"].(float64)
		if !ok {
			http.Error(response, "Invalid Claims", http.StatusUnauthorized)
			return
		}

		// Cast the userID back to int, create a new context from the old one + the userID
		intUserID := int(floatUserID)
		contextWithUserID := context.WithValue(request.Context(), UserIDKey, intUserID)
		// Pass the updated request to the handler requiring Auth
		handler(response, request.WithContext(contextWithUserID))
	}
}

func RequireAuthWS(handler http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		stringToken := request.URL.Query().Get("token")
		if stringToken == "" {
			http.Error(response, "Token Is Missing", http.StatusUnauthorized)
			return
		}

		parsedToken, err := jwt.Parse(stringToken, func(token *jwt.Token) (any, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !parsedToken.Valid {
			http.Error(response, "Invalid Or Expired Token", http.StatusUnauthorized)
			return
		}

		decodedClaims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(response, "Invalid Claims", http.StatusUnauthorized)
			return
		}

		floatUserID, ok := decodedClaims["id"].(float64)
		if !ok {
			http.Error(response, "Invalid Claims", http.StatusUnauthorized)
			return
		}

		intUserID := int(floatUserID)

		contextWithUserID := context.WithValue(request.Context(), UserIDKey, intUserID)
		handler(response, request.WithContext(contextWithUserID))
	}
} 
