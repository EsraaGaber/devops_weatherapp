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

// Connect connects to PostgreSQL instead of MySQL
func Connect(dbUser string, dbPassword string, dbHost string, dbName string, dbPort string) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	return db
}

// CreateTables creates the users table in PostgreSQL
func CreateTables(db *sql.DB) {
	cmd := `
	CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		user_name VARCHAR(50) NOT NULL UNIQUE,
		user_password CHAR(128)
	);`
	_, err := db.Exec(cmd)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func InsertUser(db *sql.DB, user User) error {
	password := md5.Sum([]byte(user.Password))
	cmd := `INSERT INTO users (user_name, user_password) VALUES ($1, $2);`
	_, err := db.Exec(cmd, user.Name, hex.EncodeToString(password[:]))
	return err
}

func GetUserByName(userName string, db *sql.DB) (User, error) {
	var user User
	cmd := `SELECT user_id, user_name, user_password FROM users WHERE user_name=$1;`
	err := db.QueryRow(cmd, userName).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil && err != sql.ErrNoRows {
		return user, err
	}
	return user, nil
}

func CreateUser(db *sql.DB, u User) (bool, error) {
	user, err := GetUserByName(u.Name, db)
	if err != nil {
		return false, err
	}
	if user != (User{}) {
		return false, nil
	} else {
		err := InsertUser(db, u)
		if err != nil {
			return false, err
		}
		return true, nil
	}
}
