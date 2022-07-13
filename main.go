package main

import (
	"fmt"
	"net/http"

	"github.com/flapan/lenslocked.com/controllers"
	"github.com/flapan/lenslocked.com/middleware"
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

	defer services.Close()
	services.AutoMigrate()
	//services.DestructiveReset()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery)

	r := mux.NewRouter()

	// Static routes
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")

	// User Routes
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	// Gallery Routes
	requiresUserMw := middleware.RequireUser{
		UserService: services.User,
	}
	r.Handle("/galleries/new", requiresUserMw.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requiresUserMw.ApplyFn(galleriesC.Create)).Methods("POST")

	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
