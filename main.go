package main

import (
	"fmt"
	"net/http"

	"github.com/flapan/lenslocked.com/controllers"
	"github.com/flapan/lenslocked.com/models"
	"github.com/gorilla/mux"
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
	services, err := models.NewServices(psqlInfo)
	must(err)
	// TODO: Fix this
	//defer us.Close()
	//us.AutoMigrate()
	//us.DestructiveReset()

	staticC := controllers.NewStatic()
	userC := controllers.NewUsers(services.User)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", userC.New).Methods("GET")
	r.HandleFunc("/signup", userC.Create).Methods("POST")
	r.Handle("/login", userC.LoginView).Methods("GET")
	r.HandleFunc("/login", userC.Login).Methods("POST")
	r.HandleFunc("/cookietest", userC.CookieTest).Methods("GET")
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
