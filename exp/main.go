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
	rows, err := db.Query(`
		SELECT * FROM users INNER JOIN orders
		ON users.id = orders.user_id`)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var userID, orderID, amount int
		var email, name, description string
		if err := rows.Scan(&userID, &name, &email, &orderID, &userID, &amount, &description); err != nil {
			panic(err)
		}
		fmt.Println("userID:", userID, "name:", name, "email:", email, "orderID:", orderID, "amount:", amount, "description:", description)
	}
	if rows.Err() != nil {
		panic(rows.Err())
	}
}
