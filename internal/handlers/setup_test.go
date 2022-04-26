package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jjang65/booking-web-app/internal/config"
	"github.com/jjang65/booking-web-app/internal/models"
	"github.com/jjang65/booking-web-app/internal/render"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	// Store Reservation type in the session
	// gob is standard library
	gob.Register(models.Reservation{})

	// Change this to ture when in production
	app.InProduction = false

	// Setup loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // Session will persist even after closing a tab
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Create templateCache initially to cache templates
	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Println(err)
		log.Fatal("cannot create template cache")
	}

	// Assign templateCache to app.TemplateCache in app config
	app.TemplateCache = tc

	// Set app.UseCache to be false, meaning no templateCache will be used
	//If set to ture, templateCache will be created, newly added temp ate won't be rendered
	// unless app server is compiled again
	app.UseCache = true

	// Passing app reference to use app config in the render package
	render.NewRenderer(&app)

	// create a new repo passing app config to be used in the handlers package
	repo := NewRepo(&app)
	// Pass pointer to repository to use in the handlers package
	NewHandlers(repo)

	mux := chi.NewRouter()

	// Recoverer middleware
	mux.Use(middleware.Recoverer)

	// NoSurf middleware for CSRF protection
	//mux.Use(NoSurf) // don't need for test; we're not testing NoSurf package

	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,                 // Only server side can access this cookie; no other client side JS can't access
		Path:     "/",                  // entire site
		Secure:   app.InProduction,     // for now, because we're using http://localhost, must be false; true is for https
		SameSite: http.SameSiteLaxMode, // SameSite Lax mode
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache creates a template cache as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	// Init map containing string key and pointer to Template
	myCache := map[string]*template.Template{}

	// Find all pages
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// Loop through all pages and if there is any template matched,
	// return parsed layouts
	for _, page := range pages {
		name := filepath.Base(page)

		// Init templateSet containing all templates
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// Find matched layouts
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		// If any matched layout, parse all layouts
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
