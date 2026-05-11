package models


import (
	"flaxat/server/db"
	"time"
)


type User struct {
	ID int					`json:"id"`
	Username string			`json:"username"`
	Email string			`json:"email"`
	Password string			`json:"password,omitempty"`
	CreatedAt time.Time		`json:"created_at"`
}

func CreateUser(username string, email string, hashedPassword string) (User, error) {
	query := `INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, created_at
	`

	var newUser User
	err := db.DB.QueryRow(query, username, email, hashedPassword).
		Scan(&newUser.ID, &newUser.Username, &newUser.Email, &newUser.CreatedAt)

	return newUser, err
}

func GetUserByID(id int) (User, error) {
	query := `SELECT id, username, email, created_at FROM users WHERE id = $1`

	var user User
	err := db.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	return user, err
}

func GetUserByEmail(email string) (User, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE email = $1`

	var user User
	err := db.DB.QueryRow(query, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)

	return user, err
}

func GetAllUsers() ([]User, error) {
	query := `SELECT id, username, email, created_at FROM users`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userList []User
	for rows.Next() {
		var current User
		err := rows.Scan(&current.ID, &current.Username, &current.Email, &current.CreatedAt)
		if err != nil {
			return nil, err
		}

		userList = append(userList, current)
	}
	return userList, nil
}
