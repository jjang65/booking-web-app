package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/jjang65/booking-web-app/internal/config"
	"github.com/jjang65/booking-web-app/internal/models"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	// Store Reservation type in the session
	// gob is standard library
	gob.Register(models.Reservation{})

	// Change this to ture when in production
	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // Session will persist even after closing a tab
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false // We're testing so should be not in https

	testApp.Session = session

	// Assign app declared in render.go to testApp created above
	app = &testApp

	os.Exit(m.Run())
}

// myWriter is a mock http.ResponseWriter
type myWriter struct {
}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {

}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
