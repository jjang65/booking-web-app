package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/jjang65/go-hello-word/pkg/config"
	"github.com/jjang65/go-hello-word/pkg/handlers"
	"github.com/jjang65/go-hello-word/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8081"

var app config.AppConfig

var session *scs.SessionManager

// main is the main application function
func main() {
	// Change this to ture when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // Session will persist even after closing a tab
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Create templateCache initially to cache templates
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	// Assign templateCache to app.TemplateCache in app config
	app.TemplateCache = tc

	// Set app.UseCache to be false, meaning no templateCache will be used
	//If set to ture, templateCache will be created, newly added temp ate won't be rendered
	// unless app server is compiled again
	app.UseCache = false

	// Passing app reference to use app config in the render package
	render.NewTemplates(&app)

	// create a new repo passing app config to be used in the handlers package
	repo := handlers.NewRepo(&app)
	// Pass pointer to repository to use in the handlers package
	handlers.NewHandlers(repo)

	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	//http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
