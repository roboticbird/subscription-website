package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"gopkg.in/robfig/cron.v3"
	"log"
	"net/http"
	"subscription-website/backend/app"
	"subscription-website/backend/db"
	"subscription-website/backend/tasks"
)

func main() {
	// connect to database
	database, err := db.DatabaseConnection()
	if err != nil {
		log.Fatal("Database connection failed: ", err.Error())
	}

	// start scheduler
	scheduler := cron.New()
	scheduler.AddFunc("@midnight", func() { tasks.OrderUpdate(database) })
	scheduler.Start()

	// start server
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
