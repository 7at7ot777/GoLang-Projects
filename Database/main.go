package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var schema = `
CREATE TABLE IF NOT EXISTS users (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   name TEXT NOT NULL,
   email TEXT NOT NULL UNIQUE,
   hashed_password BLOB NOT NULL, -- Storing as BLOB for byte slice
   created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
}

func main() {

	dbName := "users.db"
	//_ = os.Remove(dbName)
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database pinged successfully")
	//createSchema(db)
	_, err = CreateUserWithPrepared(db, "John Doe", "john@PREPARED.com", "password123")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("User created with ID: %d\n", lastId)
	//lastId, err = createUser(db, "Mideo", "john2@example.com", "password123")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("User created with ID: %d\n", lastId)
	//lastId, err = createUser(db, "Neno", "john3@example.com", "password123")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("User created with ID: %d\n", lastId)
	//lastId, err = createUser(db, "lelo", "john4@example.com", "password123")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("User created with ID: %d\n", lastId)

	//user, err := GetUserByEmail(db, "john1@example.com")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("User retrieved : %+v\n", user)

	users, _ := GetUsers(db)
	for _, user := range users {
		fmt.Printf("User retrieved : %+v\n", user)
	}

}

func createSchema(db *sql.DB) {
	_, err := db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Schema created successfully")
}

func createUser(db *sql.DB, name, email, hashedPassword string) (int64, error) {
	statement := `INSERT INTO users (name, email, hashed_password) VALUES (?, ?, ?)`
	hp, err := bcrypt.GenerateFromPassword([]byte(hashedPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	result, err := db.Exec(statement, name, email, string(hp))
	if err != nil {
		return 0, err
	}
	fmt.Println("User created successfully")
	return result.LastInsertId()
}

// prepared: compiled and saved query, execution, and optimization plan which saves resources and improves performance, and also workss against sql injection.
func CreateUserWithPrepared(db *sql.DB, name, email, hashedPassword string) (int64, error) {
	statement, err := db.Prepare(`INSERT INTO users (name, email, hashed_password) VALUES (?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer statement.Close()
	hp, err := bcrypt.GenerateFromPassword([]byte(hashedPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	result, err := statement.Exec(name, email, string(hp))
	if err != nil {
		return 0, err
	}
	fmt.Println("User created successfully")
	return result.LastInsertId()
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	statement := `SELECT * FROM users WHERE email = ?`
	row := db.QueryRow(statement, email)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func GetUsers(db *sql.DB) ([]User, error) {
	statement := `SELECT * from users`
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users, nil
}
