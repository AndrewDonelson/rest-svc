package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AndrewDonelson/golog"
	"github.com/AndrewDonelson/rest-svc/db"
	"github.com/gorilla/mux"
)

// ServerApp defines the reguired functions that need to be implemented
type ServerAppInterface interface {
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

type ServerAppState int

const (
	AppStateUnknown      ServerAppState = 0
	AppStateInitializing ServerAppState = 1
	AppStateInitialized  ServerAppState = 2
	AppStateOnline       ServerAppState = 3
	AppStateOffline      ServerAppState = 4
	AppStateShuttingDown ServerAppState = 5
	AppStateShutdown     ServerAppState = 6
	AppStateExiting      ServerAppState = 7
)

type RouterShutdownFunc func(r *mux.Router)
type DatabaseShutdownFunc func(d *sql.DB)
type NetworkShutdownFunc func(hostaddresses []string)
type AppShutdownFunc func(notifyEmail string)

type ShutdownMethods interface {
	NiceRouterShudown(r *mux.Router)
	NiceDatbaseShutdown(d *sql.DB)
	NiceNetworkShutdown(hostaddresses []string)
	NiceAppShutdown(notifyEmail string)
}

type ShutdownFuncs struct {
	RouterFunc  RouterShutdownFunc
	Database    DatabaseShutdownFunc
	Network     NetworkShutdownFunc
	Application AppShutdownFunc
}

type ServerConfig struct {
	Name    string
	Version string
}

// App is the main application object and must implement all methods of ServerApp
type ServerApp struct {
	state         ServerAppState
	config        ServerConfig
	Router        *mux.Router
	Database      *sql.DB
	ShutdownFuncs ShutdownFuncs
}

var Svr *ServerApp

func init() {

	Svr = new(ServerApp)
	Svr.config = ServerConfig{Name: "rest-svc", Version: "1.0.0"}
	golog.Log.Options = golog.Options{Module: Svr.config.Name, Environment: golog.EnvDevelopment, SmartError: true}
	golog.Log.Info("ServerApp Initialized")

}

// InitRouter called by Startup to initialize all routes handled by app
func (app *ServerApp) InitRouter() {
	golog.Log.Info("Initializing router")
	app.Router = mux.NewRouter().StrictSlash(true)
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
func (app *ServerApp) InitDatabase() {
	var err error
	app.Database, err = db.CreateDatabase()
	if err != nil {
		golog.Log.Warningf("Database connection failed: %s", err.Error())
		golog.Log.Info("Database access will be disabled")
	}
}

// Router Handler Functions (local)

func (app *ServerApp) handlerIndexFunction(w http.ResponseWriter, r *http.Request) {
	golog.Log.HandlerLog(w, r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s v%s", app.config.Name, app.config.Version)
}

func (app *ServerApp) handlerGetFunction(w http.ResponseWriter, r *http.Request) {
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

func (app *ServerApp) handlerPostFunction(w http.ResponseWriter, r *http.Request) {
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

func (app *ServerApp) Main() {
	var err error

	app.InitDatabase()
	app.InitRouter()

	golog.Log.Info("Listening on port 8080")
	err = http.ListenAndServe(":8080", app.Router)
	golog.Log.Fatal(err.Error())

}
