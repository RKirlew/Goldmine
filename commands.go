package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	_ "github.com/mattn/go-sqlite3"
)

func add(db *sql.DB, err error, email string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter site: ")
	site, _ := reader.ReadString('\n')

	fmt.Print("Enter password: ")
	password, _ := terminal.ReadPassword(int(os.Stdin.Fd()))

	createdAt := time.Now().Format("2006-01-02 15:04:05")

	_, err = db.Exec("INSERT INTO passwords (site, email, password, created_at) VALUES (?, ?, ?, ?)",
		strings.TrimSpace(site), email, string(password), createdAt)

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Password entry added successfully!")
	return nil
}
