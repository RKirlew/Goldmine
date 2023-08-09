package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func getUserByUsername(email string) (string, error) {
	// Query the database to retrieve the user's hashed password by their email
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		return "", fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	query := "SELECT password FROM 'users' WHERE email =?"
	fmt.Println("Query:", query) // Print the complete query
	row := db.QueryRow(query, email)

	var storedPassword string
	err = row.Scan(&storedPassword)
	//fmt.Println(storedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			// User not found
			fmt.Println(err)
			return "grw3", nil
		}
		return "", fmt.Errorf("error fetching data: %v", err)
	}

	return storedPassword, nil
}

func Login(username, password string) (*User, error) {
	// Retrieve the user's hashed password from the database based on the username
	hashedPassword, err := getUserByUsername(username)
	if err != nil {
		return nil, err
	}

	// Verify the password
	if !verifyPassword(password, hashedPassword) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Password is correct, return the user
	user := &User{
		Username:     username,
		PasswordHash: hashedPassword,
	}

	return user, nil
}

func verifyPassword(password, passwordHash string) bool {
	// Hash the provided password and compare it with the stored password hash
	fmt.Println("pass:", password)
	fmt.Println("pasHasH", passwordHash)
	fmt.Println(hashPassword(passwordHash))
	return hashPassword(password) == passwordHash
}
