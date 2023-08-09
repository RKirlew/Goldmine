// register.go

package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func Register(username, email string, password string) (*User, error) {
	// Hash the password using the utility function from user.go
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		fmt.Println("Error opening database:", err)

	} else {
		defer db.Close()
		passwordHash := hashPassword(strings.TrimSpace(password))
		currentTime := time.Now()
		// Create a new User instance
		user := &User{
			Username:     username,
			PasswordHash: passwordHash,
		}

		// Store the user in your database or data store

		// insert
		stmt, err := db.Prepare("INSERT INTO users(username,email,password,created_at) values(?,?,?,?)")
		checkErr(err)

		res, err := stmt.Exec(username, strings.TrimSpace(email), passwordHash, currentTime)
		checkPassErr(err)

		id, err := res.LastInsertId()
		fmt.Println(id)

		//checkErr(err)
		// For simplicity, we'll just return the user instance here.

		return user, err
	}
	return nil, nil
}
