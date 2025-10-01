package authdb

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"

	_ "github.com/lib/pq"
)

type User struct {
	ID       int    `json:"user_id"`
	Name     string `json:"user_name"`
	Password string `json:"user_password"`
}

// Connect to PostgreSQL database
func Connect(dbUser, dbPassword, dbHost, dbName string) *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s port=5432 user=%s password=%s dbname=%s sslmode=require",
		dbHost, dbUser, dbPassword, dbName,
	)

    fmt.Println("🔌 Connecting to Postgres with:")
    fmt.Println("  Host:", dbHost)
    fmt.Println("  User:", dbUser)
    fmt.Println("  DB:", dbName)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error connecting to DB:", err)
	}
	return db
}

// CreateTables creates the users table if it doesn't exist
func CreateTables(db *sql.DB) {
	cmd := `
	CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		user_name VARCHAR(50) NOT NULL UNIQUE,
		user_password CHAR(32) NOT NULL
	);`
	_, err := db.Exec(cmd)
	if err != nil {
		fmt.Println("Error creating table:", err)
	}
}

// InsertUser inserts a new user into the users table
func InsertUser(db *sql.DB, user User) error {
	password := md5.Sum([]byte(user.Password))
	cmd := `
	INSERT INTO users (user_name, user_password)
	VALUES ($1, $2);`
	_, err := db.Exec(cmd, user.Name, hex.EncodeToString(password[:]))
	return err
}

// GetUserByName retrieves a user by username
func GetUserByName(db *sql.DB, userName string) (User, error) {
	var user User
	cmd := `SELECT user_id, user_name, user_password FROM users WHERE user_name=$1`
	row := db.QueryRow(cmd, userName)
	err := row.Scan(&user.ID, &user.Name, &user.Password)
	if err == sql.ErrNoRows {
		return User{}, nil // explicitly no user found
	}
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// CreateUser creates a new user if not exists
func CreateUser(db *sql.DB, u User) (bool, error) {
	existingUser, err := GetUserByName(db, u.Name)
	if err != nil {
		return false, err
	}
	if existingUser != (User{}) {
		return false, nil // user already exists
	}
	err = InsertUser(db, u)
	if err != nil {
		return false, err
	}
	return true, nil
}
