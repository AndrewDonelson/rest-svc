package main

import (
	"net/http"

	"github.com/AndrewDonelson/golog"
	"github.com/AndrewDonelson/rest-svc/app"
	"github.com/AndrewDonelson/rest-svc/db"
	"github.com/gorilla/mux"
)

//var Log *golog.Logger

func main() {
	golog.Log.Options = golog.Options{Module: "rest-svc", Environment: golog.EnvDevelopment, SmartError: true}

	database, err := db.CreateDatabase()
	if err != nil {
		golog.Log.Warningf("Database connection failed: %s", err.Error())
		golog.Log.Info("Database access will be disabled")
	}

	app := &app.App{
		Router:   mux.NewRouter().StrictSlash(true),
		Database: database,
	}

	app.InitRouter()

	golog.Log.Info("Listening on port 8080")
	err = http.ListenAndServe(":8080", app.Router)
	golog.Log.Fatal(err.Error())
}
