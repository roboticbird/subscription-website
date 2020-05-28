package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"subscription-website/backend/app"
	"subscription-website/backend/db"
)

func main() {
	database, err := db.DatabaseConnection()
	if err != nil {
		log.Fatal("Database connection failed: ", err.Error())
	}

	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	app := &app.App{
		Router:   mux.NewRouter().StrictSlash(true),
		Store:    sessions.NewCookieStore(authKeyOne, encryptionKeyOne),
		Database: database,
	}

	app.Setup()

	log.Fatal(http.ListenAndServe(":8080", app.Router))
}
