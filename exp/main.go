package main

import (
	"fmt"

	"github.com/flapan/lenslocked.com/models"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host   = "localhost"
	port   = 5432
	dbuser = "mikkel"
	dbname = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, dbuser, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	//us.DestructiveReset()
	user, err := us.ByID(2)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

}
