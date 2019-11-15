package app

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/AndrewDonelson/golog"
	"github.com/AndrewDonelson/rest-svc/db"
	"github.com/gorilla/mux"
)

// ServerApp defines the reguired functions that need to be implemented
type ServerApp interface {
	// Startup calls Init* methods and displays standard app header
	Startup()
	// InitEnvironment called by Startup to intiialize apps environment variables
	InitEnvironment()
	// InitConfig called by Startup to load the apps configuration
	InitConfig()
	// InitFlags called by Startup to process the apps flags which can be used to override
	// Environment and Config
	InitFlags()
	// InitDatabase called by Startup to initialize the apps database (if any)
	InitDatabase()
	// InitRouter called by Startup to initialize all routes handled by app
	InitRouter()
	// Listen Invoke after Startup with no errors to handle the main loop as well as
	// a graceful shutdown
	Listen()
	// Shutdown is called from Listen does not need to be called directly
	Shutdown()
}

// App is the main application object and must implement all methods of ServerApp
type App struct {
	Router   *mux.Router
	Database *sql.DB
}

// InitRouter called by Startup to initialize all routes handled by app
func (app *App) InitRouter() {
	golog.Log.Info("Initializing router")
	app.Router.
		Methods("GET").
		Path("/").
		HandlerFunc(app.handlerIndexFunction)

	app.Router.
		Methods("GET").
		Path("/endpoint/{id}").
		HandlerFunc(app.handlerGetFunction)

	app.Router.
		Methods("POST").
		Path("/endpoint").
		HandlerFunc(app.handlerPostFunction)
}

// InitDatabase called by Startup to initialize the apps database (if any)
func (app *App) InitDatabase() (err error) {
	app.Database, err = db.CreateDatabase()
	if err != nil {
		golog.Log.Warningf("Database connection failed: %s", err.Error())
		golog.Log.Info("Database access will be disabled")
	}

	return
}

// Router Handler Functions (local)

func (app *App) handlerIndexFunction(w http.ResponseWriter, r *http.Request) {
	golog.Log.HandlerLog(w, r)

}

func (app *App) handlerGetFunction(w http.ResponseWriter, r *http.Request) {
	golog.Log.HandlerLog(w, r)
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		golog.Log.Fatal("No ID in the path")
	}

	dbdata := &DbData{}
	if app.Database != nil {
		err := app.Database.QueryRow("SELECT id, date, name FROM `test` WHERE id = ?", id).Scan(&dbdata.ID, &dbdata.Date, &dbdata.Name)
		if err != nil {
			golog.Log.Fatal("Database SELECT failed")
		}
	}

	golog.Log.Notice("You fetched a thing!")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dbdata); err != nil {
		panic(err)
	}
}

func (app *App) handlerPostFunction(w http.ResponseWriter, r *http.Request) {
	golog.Log.HandlerLog(w, r)
	if app.Database != nil {
		_, err := app.Database.Exec("INSERT INTO `test` (name) VALUES ('myname')")
		if err != nil {
			golog.Log.Fatal("Database INSERT failed")
		}
	}
	golog.Log.Notice("You called a thing!")
	w.WriteHeader(http.StatusOK)
}
