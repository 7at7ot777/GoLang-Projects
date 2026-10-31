package main

import (
	"database/sql"
	"firstProgram/RepositoryPattern/repository"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	dbName := "users.db"
	db, err := connectToDatabase(dbName)
	checkError(err)
	defer db.Close()

	fmt.Println("database connection established")

	repo := repository.NewSQLUserRepository(db)
	user, err := repo.CreateUser("John Doe", "johdfghjkn@example.com", "nenene", "asdadasda")
	checkError(err)
	fmt.Println(user)
}
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func connectToDatabase(name string) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
