package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/jjang65/booking-web-app/internal/config"
	"github.com/jjang65/booking-web-app/internal/driver"
	"github.com/jjang65/booking-web-app/internal/handlers"
	"github.com/jjang65/booking-web-app/internal/helpers"
	"github.com/jjang65/booking-web-app/internal/models"
	"github.com/jjang65/booking-web-app/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8081"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	//http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// Store Reservation type in the session
	// gob is standard library
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// Change this to ture when in production
	app.InProduction = false

	// Setup loggers
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // Session will persist even after closing a tab
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Connect to db
	log.Println("connecting to db")
	db, err := driver.ConnectSQL("host=172.18.0.3 port=5432 dbname=bookings user=root password=root")
	if err != nil {
		log.Fatal("cannot connect to db")
	}
	log.Println("connected to db")

	// Create templateCache initially to cache templates
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Println(err)
		log.Fatal("cannot create template cache")
		return nil, err
	}

	// Assign templateCache to app.TemplateCache in app config
	app.TemplateCache = tc

	// Set app.UseCache to be false, meaning no templateCache will be used
	//If set to ture, templateCache will be created, newly added temp ate won't be rendered
	// unless app server is compiled again
	app.UseCache = false

	// Passing app reference to use app config in the render package
	render.NewRenderer(&app)

	// Passing app reference to helpers
	helpers.NewHelpers(&app)

	// create a new repo passing app config to be used in the handlers package
	repo := handlers.NewRepo(&app, db)
	// Pass pointer to repository to use in the handlers package
	handlers.NewHandlers(repo)

	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)
	return db, nil
}
