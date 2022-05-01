package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jjang65/booking-web-app/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{
		"home",
		"/",
		"GET",
		http.StatusOK,
	},
	{
		"about",
		"/about",
		"GET",
		http.StatusOK,
	},
	{
		"gq",
		"/generals-quarters",
		"GET",
		http.StatusOK,
	},
	{
		"ms",
		"/majors-suite",
		"GET",
		http.StatusOK,
	},
	{
		"sa",
		"/search-availability",
		"GET",
		http.StatusOK,
	},
	{
		"contact",
		"/contact",
		"GET",
		http.StatusOK,
	},
	//{
	//	"post-search-avail",
	//	"/search-availability",
	//	"POST",
	//	[]postData{
	//		{key: "start", value: "2022-01-01"},
	//		{key: "end", value: "2022-01-31"},
	//	},
	//	http.StatusOK,
	//},
	//{
	//	"post-search-avail-json",
	//	"/search-availability-json",
	//	"POST",
	//	[]postData{
	//		{key: "start", value: "2022-01-01"},
	//		{key: "end", value: "2022-01-31"},
	//	},
	//	http.StatusOK,
	//},
	//{
	//	"post-make-reservation",
	//	"/make-reservation",
	//	"POST",
	//	[]postData{
	//		{key: "first_name", value: "John"},
	//		{key: "last_name", value: "Smith"},
	//		{key: "email", value: "me@here.com"},
	//		{key: "phone", value: "555-555-5555"},
	//	},
	//	http.StatusOK,
	//},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}

}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room:   models.Room{ID: 1, RoomName: "General's Quarters"},
	}
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	// Get Context containing Session
	ctx := getCtx(req)
	// Now request has context with session
	req = req.WithContext(ctx)

	// Get new ResponseRecorder that is basically a mock http response
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	//	test case where reservation is not in session (reset everything)
	// In this case, http status code should be 307
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	log.Println("rr.Code: ", rr.Code)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//	test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	//reservation.RoomID = 100
	//session.Put(ctx, "reservation", reservation)
	//
	handler.ServeHTTP(rr, req)
	log.Println("rr.Code: ", rr.Code)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	// Test valid case
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=j@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=j@123123123")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	// Get Context containing Session
	ctx := getCtx(req)
	// Now request has context with session
	req = req.WithContext(ctx)

	// Set header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Get new ResponseRecorder that is basically a mock http response
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	// http status code should be StatusSeeOther (303)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	//	Test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	// Get Context containing Session
	ctx = getCtx(req)
	// Now request has context with session
	req = req.WithContext(ctx)

	// Set header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Get new ResponseRecorder that is basically a mock http response
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	// http status code should be StatusTemporaryRedirect (307)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test for invalid start date after resetting request.body
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=j@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=j@123123123")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	// Get Context containing Session
	ctx = getCtx(req)
	// Now request has context with session
	req = req.WithContext(ctx)

	// Set header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Get new ResponseRecorder that is basically a mock http response
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	// http status code should be StatusTemporaryRedirect (307)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test for invalid end date after resetting request.body
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=j@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=j@123123123")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	// Get Context containing Session
	ctx = getCtx(req)
	// Now request has context with session
	req = req.WithContext(ctx)

	// Set header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Get new ResponseRecorder that is basically a mock http response
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	// http status code should be StatusTemporaryRedirect (307)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong response code for invalid end date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test for invalid room id
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=j@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=j@123123123")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	// Get Context containing Session
	ctx = getCtx(req)
	// Now request has context with session
	req = req.WithContext(ctx)

	// Set header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Get new ResponseRecorder that is basically a mock http response
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	// http status code should be StatusTemporaryRedirect (307)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong response code for invalid room id: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test for invalid data
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=J")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=j@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=j@123123123")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	// Get Context containing Session
	ctx = getCtx(req)
	// Now request has context with session
	req = req.WithContext(ctx)

	// Set header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Get new ResponseRecorder that is basically a mock http response
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	// http status code should be StatusTemporaryRedirect (307)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("reservation handler returned wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// Test for failure to insert restriction into database
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=j@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=j@123123123")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000") // should fail expectedly

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	// Get Context containing Session
	ctx = getCtx(req)
	// Now request has context with session
	req = req.WithContext(ctx)

	// Set header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Get new ResponseRecorder that is basically a mock http response
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	// http status code should be StatusTemporaryRedirect (307)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler failed when trying to fail inserting resrvation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test for failure to insert reservation into database
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=j@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=j@123123123")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2") // Should fail as expected

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	// Get Context containing Session
	ctx = getCtx(req)
	// Now request has context with session
	req = req.WithContext(ctx)

	// Set header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Get new ResponseRecorder that is basically a mock http response
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	// http status code should be StatusTemporaryRedirect (307)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler failed when trying to fail inserting resrvation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// Test rooms are not available
	postedData := url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("room_id", "1")

	// create request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(postedData.Encode()))

	// get context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "x-www-form-urlencoded")

	// make handlerfunc
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	// get response recorder
	rr := httptest.NewRecorder()

	// make request to our handler
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Errorf("failed to parse json")
	}

	// http status code should be StatusTemporaryRedirect (200)
	if rr.Code != http.StatusOK {
		t.Errorf("AvailabilityJson handler failed: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println("getCtx::err", err)
	}
	return ctx
}
