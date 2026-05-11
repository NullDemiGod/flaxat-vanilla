package handlers

import (
	"database/sql"
	"encoding/json"
	"flaxat/server/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Helper function that json-encodes custom structs and writes them to the response body, also sets status codes
func writeJSON(response http.ResponseWriter, status int, data any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	json.NewEncoder(response).Encode(data)
}

// Helper function that generates JWT Token from a user ID
func createJWToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"id": userID,
		"exp": time.Now().Add(2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// The register handler, the interface between the frontend and database, goes for all the handlers
func Register(response http.ResponseWriter, request *http.Request) {
	// Define an inline struct that defines the structure of the request body
	var bodyStructure struct {
		Username string		`json:"username"`
		Email string		`json:"email"`
		Password string		`json:"password"`
	}

	// Decode the incoming request body into the bodyStructure object
	err := json.NewDecoder(request.Body).Decode(&bodyStructure)
	if err != nil {
		writeJSON(response, http.StatusBadRequest, map[string]string{
			"error": "Invalid Request Body",
		})
		return
	}

	// Enforce NOT NULL CONSTRAINTS, all fields are required
	if bodyStructure.Username == "" || bodyStructure.Email == "" || bodyStructure.Password == "" {
		writeJSON(response, http.StatusBadRequest, map[string]string{
			"error": "All Fields Are Required",
		})
		return
	}

	// Hash the password using bcrypt, 12 is the number of rounds of encryption
	// Strong enough to resist brute force decryption
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(bodyStructure.Password), 12)
	if err != nil {
		writeJSON(response, http.StatusInternalServerError, map[string]string{
			"error": "Failed To Hash Password",
		})
		return
	}

	createdUser, err := models.CreateUser(bodyStructure.Username, bodyStructure.Email, string(hashedPassword))
	if err != nil {
		writeJSON(response, http.StatusBadRequest, map[string]string{
			"error": "Username Or Email Already Exists",
		})
		return
	}

	token, err := createJWToken(createdUser.ID)
	if err != nil {
		writeJSON(response, http.StatusInternalServerError, map[string]string{
			"error": "Failed To Generate JWT Token",
		})
		return
	}
	
	writeJSON(response, http.StatusCreated, map[string]any{
		"user": createdUser,
		"token": token,
	})
}

func Login(response http.ResponseWriter, request *http.Request) {
	var bodyStructure struct {
		Email string 		`json:"email"`
		Password string 	`json:"password"`
	}

	err := json.NewDecoder(request.Body).Decode(&bodyStructure)
	if err != nil {
		writeJSON(response, http.StatusBadRequest, map[string]string{
			"error": "Invalid Request Body",
		})
		return
	}

	if bodyStructure.Email == "" || bodyStructure.Password == "" {
		writeJSON(response, http.StatusBadRequest, map[string]string{
			"error": "All Fields Are Required",
		})
		return
	}

	user, err := models.GetUserByEmail(bodyStructure.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(response, http.StatusBadRequest, map[string]string{
				"error": "Invalid Credentials",
			})
		} else {
			writeJSON(response, http.StatusInternalServerError, map[string]string{
				"error": "Something Went Wrong",
			})
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(bodyStructure.Password))
	if err != nil {
		writeJSON(response, http.StatusBadRequest, map[string]string{
			"error": "Invalid Credentials",
		})
		return
	}

	token, err := createJWToken(user.ID)
	if err != nil {
		writeJSON(response, http.StatusInternalServerError, map[string]string{
			"error": "Failed To Generate JWT Token",
		})
		return
	}

	user.Password = ""
	writeJSON(response, http.StatusOK, map[string]any{
		"user": user,
		"token": token,
	})
}

func GetUser(response http.ResponseWriter, request *http.Request) {
	userIDString := request.PathValue("userID")

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		writeJSON(response, http.StatusBadRequest, map[string]string{
			"error": "Invalid User ID",
		})
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(response, http.StatusBadRequest, map[string]string{
				"error": "User doesn't exist",
			})
		} else {
			writeJSON(response, http.StatusInternalServerError, map[string]string{
				"error": "Something Went Wrong",
			})
		}
		return
	}

	writeJSON(response, http.StatusOK, user)
}

func GetAllUsers(response http.ResponseWriter, request *http.Request) {
	userList, err := models.GetAllUsers()
	if err != nil {
		writeJSON(response, http.StatusInternalServerError, map[string]string{
			"error": "Failed To Fetch Users",
		})
		return
	}

	writeJSON(response, http.StatusOK, userList)
}
