package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "mikkel"
	dbname = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var id int

	err = db.QueryRow(`
	INSERT INTO users(name, email)
	VALUES($1, $2)
	RETURNING id`,
		"Hugo Gadegaard Schiller", "hugogs@gmail.com").Scan(&id)

	if err != nil {
		panic(err)
	}
	fmt.Println("id is...", id)
}
