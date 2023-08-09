package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type PasswordEntry struct {
	Site      string
	Email     string
	Password  string
	CreatedAt time.Time
}

func contains(s string, str string) bool {
	return strings.Contains(s, str)
}
func checkErr(err error) {
	if err != nil {
		fmt.Println("Error opening database")
		//panic(err)
	}

}
func checkPassErr(err error) {
	if err != nil {
		fmt.Println("Error inserting password")
		//panic(err)
	}

}
func generateSecure() string {
	const length = 16
	fmt.Println("Please note: this password wont be stored since you are not logged in")

	// Create a byte slice to store the random data
	var tags []string

	tags = make([]string, 3)

	tags[0] = "-g"
	tags[1] = "-t"
	tags[2] = "-b"
	randomBytes := make([]byte, length)

	// Generate secure random data
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println("Error generating secure random data:", err)
		return ""
	}

	// Encode the random data to base64 to make it human-readable
	secureRandomString := base64.URLEncoding.EncodeToString(randomBytes)

	// Print the generated secure random string
	fmt.Println("Generated Secure Random String:", secureRandomString)
	return secureRandomString

}

func main() {

	//user fields
	var userName string
	var newUser User
	_ = newUser
	var emailField string

	var masterPass string

	userPtr := flag.String("email", "tester@test.com", "the default email")
	_ = userPtr

	inputtedPassPtr := flag.String("password", "testerpass123", "the default password")
	_ = inputtedPassPtr

	loginPtr := flag.Bool("l", false, "a bool for determinining if user is logging in or not")
	_ = loginPtr
	db, err := sql.Open("sqlite3", "./foo.db")
	_ = db
	checkErr(err)

	createTableSQL := `
				CREATE TABLE IF NOT EXISTS passwords (
					site VARCHAR(255) NOT NULL,
					email VARCHAR(255) NOT NULL,
					password VARCHAR(255) NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				)
			`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createUserTable := `
				CREATE TABLE IF NOT EXISTS users (
					username VARCHAR(255) NOT NULL,
					email VARCHAR(255) NOT NULL,
					password VARCHAR(255) NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				)
			`
	_, err = db.Exec(createUserTable)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("Tables created successfully!")

	cipherKey := []byte("asuperstrong32bitpasswordgohere!") //32 bit key for AES-256
	//cipherKey := []byte("asuperstrong24bitpasswor") //24 bit key for AES-192
	//cipherKey := []byte("asuperstrong16bi") //16 bit key for AES-128

	reader := bufio.NewReader(os.Stdin)

	var message string
	if len(os.Args) == 2 && contains(os.Args[1], "-g") {

		generateSecure()
	} else if len(os.Args) > 2 && contains(os.Args[1], "-g") {
		fmt.Println("Usage: go run main.go -g")
	}

	//IF no command line argument is given:

	if len(os.Args) == 1 {

		//User setup
		fmt.Printf("\n\tNo command line argument found. Assuming account setup\n")
		fmt.Printf("\n\tEnter a username to register: ")
		userName, _ = reader.ReadString('\n')
		fmt.Printf("\tEnter email for %s ", userName)
		emailField, _ = reader.ReadString('\n')
		fmt.Printf("\tEnter masterpass for %s ", userName)
		masterPass, _ = reader.ReadString('\n')
		fmt.Println(emailField, userName, masterPass)

		newUser, err := Register(userName, emailField, masterPass)
		fmt.Println(newUser)

		if err != nil {
			fmt.Println("Error registering:", err, newUser)
			return
		} else {

			fmt.Printf("Successfully registered:%s", newUser)

		}
	}
	if len(os.Args) > 1 && contains(os.Args[1], "-l") {
		fmt.Println("Initializing Login...")

		flag.Parse()
		loggedInUser, err := Login(*userPtr, *inputtedPassPtr)
		if err != nil {
			fmt.Println("Error logging in:", err)
			return
		}

		fmt.Printf("Logged In User: %s\n", loggedInUser.Username)

		for {
			fmt.Print("Enter command > ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}

			// Remove leading/trailing whitespaces and convert to lowercase for case-insensitive comparison
			command := strings.TrimSpace(strings.ToLower(input))

			switch command {
			case "1":
				// View Passwords - Implement this function
				//viewBubbleUi()
				time.Sleep(3)
				cmd := exec.Command("cmd", "/c", "cls")
				cmd.Stdout = os.Stdout
				cmd.Run()
			case "q":
				os.Exit(1)
			case "add":
				add(db, err, loggedInUser.Username)
			case "help":
				showHelpScreen()
			case "generate":
				generateSecure()
			case "all":
				fetchAll(loggedInUser.Username)

			default:
				fmt.Println("\nInvalid command. Please try again.")
			}
		}
	}

	if len(os.Args) > 1 && strings.Contains(os.Args[1], "-e") { //Make the message equal to the command line argument
		message = os.Args[2]
		fmt.Println(message)
		encrypted, err := encrypt(cipherKey, message)

		//IF the encryption failed:
		if err != nil {
			//Print error message:
			log.Println(err)
			os.Exit(-2)
		}

		//Print the key and cipher text:
		fmt.Printf("\n\tCIPHER KEY: %s\n", string(cipherKey))
		fmt.Printf("\tENCRYPTED: %s\n", encrypted)

		//Decrypt the text:
		decrypted, err := decrypt(cipherKey, encrypted)

		//IF the decryption failed:
		if err != nil {
			log.Println(err)
			os.Exit(-3)
		}

		//Print re-decrypted text:
		fmt.Printf("\tDECRYPTED: %s\n\n", decrypted)
	}

}

func fetchAll(user string) {
	//fmt.Println("More:", user)
	query := "SELECT site, email, password, created_at FROM passwords WHERE email=?"
	db, err := sql.Open("sqlite3", "./foo.db")
	_ = db
	// Execute the query
	rows, err := db.Query(query, user)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate through the rows and populate PasswordEntry structs
	var entries []PasswordEntry
	for rows.Next() {
		var entry PasswordEntry
		err := rows.Scan(&entry.Site, &entry.Email, &entry.Password, &entry.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		entries = append(entries, entry)
	}

	// Print the fetched entries
	for _, entry := range entries {
		fmt.Printf("Site: %s\nEmail: %s\nPassword: %s\nCreated At: %s\n\n",
			entry.Site, entry.Email, entry.Password, entry.CreatedAt)
	}
}
func encrypt(key []byte, message string) (encoded string, err error) {
	//Create byte array from the input string
	plainText := []byte(message)

	//Create a new AES cipher using the key
	block, err := aes.NewCipher(key)

	//IF NewCipher failed, exit:
	if err != nil {
		return
	}

	//Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	//iv is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	//Encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//Return string encoded in base64
	return base64.RawStdEncoding.EncodeToString(cipherText), err
}

func decrypt(key []byte, secure string) (decoded string, err error) {
	//Remove base64 encoding:
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)

	//IF DecodeString failed, exit:
	if err != nil {
		return
	}

	//Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher(key)

	//IF NewCipher failed, exit:
	if err != nil {
		return
	}

	//IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	//Decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), err
}
