package main

import (
	"fmt"
	"github.com/jjang65/booking-web-app/internal/helpers"
	"github.com/justinas/nosurf"
	"net/http"
)

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the page")
		next.ServeHTTP(w, r)
	})
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

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(
		// create own handler func where to call IsAuthenticate()
		// then call next.ServeHTTP() to continue
		func(w http.ResponseWriter, r *http.Request) {
			if !helpers.IsAuthenticated(r) {
				session.Put(r.Context(), "error", "Please login")
				http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			}
			next.ServeHTTP(w, r)
		})
}
