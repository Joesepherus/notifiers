package controllers_test

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tradingalerts/controllers"
	database "tradingalerts/db"
	"tradingalerts/services/alertsService"
	"tradingalerts/services/loggingService"
	"tradingalerts/services/userService"
	"tradingalerts/templates"
)

var db *sql.DB

func TestMain(m *testing.M) {
	// Initialize templates and database
	log.Print("Initializing tests")
	templates.InitTemplates("../templates")
	db = database.InitDB("../alerts_test.db")
	defer db.Close() // Ensure database connection is closed after tests

	// Pass the db connection to services
	alertsService.SetDB(db)
	userService.SetDB(db)
	loggingService.SetDB(db)

	// Run the tests
	m.Run()
}

func TestPageHandler_Home(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.PageHandler)
	handler.ServeHTTP(rr, req)

	// Check the response code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the title in the response body
	body := rr.Body.String()
	expectedTitle := "Trading Alerts"
	if !strings.Contains(body, `<title>`+expectedTitle+`</title>`) {
		t.Errorf("handler returned unexpected title: got %v", body)
	}
}

func TestPageHandler_Pricing(t *testing.T) {
	req, err := http.NewRequest("GET", "/pricing", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.PageHandler)
	handler.ServeHTTP(rr, req)

	// Check the response code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the title in the response body
	body := rr.Body.String()
	expectedTitle := "Pricing - Trading Alerts"
	if !strings.Contains(body, `<title>`+expectedTitle+`</title>`) {
		t.Errorf("handler returned unexpected title: got %v", body)
	}
}

func TestPageHandler_About(t *testing.T) {
	req, err := http.NewRequest("GET", "/about", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.PageHandler)
	handler.ServeHTTP(rr, req)

	// Check the response code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the title in the response body
	body := rr.Body.String()
	expectedTitle := "About - Trading Alerts"
	if !strings.Contains(body, `<title>`+expectedTitle+`</title>`) {
		t.Errorf("handler returned unexpected title: got %v", body)
	}
}

// func TestPageHandler_Alerts(t *testing.T) {
// 	log.Print("Running Alerts Page Handler Test")
// 	req, err := http.NewRequest("GET", "/alerts", nil)
// 	if err != nil {
// 		t.Fatalf("Could not create request: %v", err)
// 	}

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(controllers.PageHandler)
// 	handler.ServeHTTP(rr, req)

// 	// Check the response code
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	// Check the title in the response body
// 	body := rr.Body.String()
// 	expectedTitle := "Alerts - Trading Alerts"
// 	if !strings.Contains(body, `<title>`+expectedTitle+`</title>`) {
// 		t.Errorf("handler returned unexpected title: got %v", body)
// 	}
// }

func TestPageHandler_Profile(t *testing.T) {
	req, err := http.NewRequest("GET", "/profile", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.PageHandler)
	handler.ServeHTTP(rr, req)

	// Check the response code
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	expectedLocation := "/error?message=You+need+to+be+logged+in"
	if rr.Header().Get("Location") != expectedLocation {
		t.Errorf("handler returned wrong redirect location: got %v want %v", rr.Header().Get("Location"), expectedLocation)
	}
}

func TestPageHandler_ResetPasswordSent(t *testing.T) {
	req, err := http.NewRequest("GET", "/reset-password-sent", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.PageHandler)
	handler.ServeHTTP(rr, req)

	// Check the response code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the title in the response body
	body := rr.Body.String()
	expectedTitle := "Reset Password - Trading Alerts"
	if !strings.Contains(body, `<title>`+expectedTitle+`</title>`) {
		t.Errorf("handler returned unexpected title: got %v", body)
	}
}

func TestPageHandler_ResetPasswordSuccess(t *testing.T) {
	req, err := http.NewRequest("GET", "/reset-password-success", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.PageHandler)
	handler.ServeHTTP(rr, req)

	// Check the response code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the title in the response body
	body := rr.Body.String()
	expectedTitle := "Reset Password Success - Trading Alerts"
	if !strings.Contains(body, `<title>`+expectedTitle+`</title>`) {
		t.Errorf("handler returned unexpected title: got %v", body)
	}
}

func TestPageHandler_SubscriptionSuccessful(t *testing.T) {
	req, err := http.NewRequest("GET", "/subscription-success", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.PageHandler)
	handler.ServeHTTP(rr, req)

	// Check the response code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the title in the response body
	body := rr.Body.String()
	expectedTitle := "Subscription Successful - Trading Alerts"
	if !strings.Contains(body, `<title>`+expectedTitle+`</title>`) {
		t.Errorf("handler returned unexpected title: got %v", body)
	}
}

func TestPageHandler_SubscriptionCancelled(t *testing.T) {
	req, err := http.NewRequest("GET", "/subscription-cancel", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.PageHandler)
	handler.ServeHTTP(rr, req)

	// Check the response code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the title in the response body
	body := rr.Body.String()
	expectedTitle := "Subscription Cancelled - Trading Alerts"
	if !strings.Contains(body, `<title>`+expectedTitle+`</title>`) {
		t.Errorf("handler returned unexpected title: got %v", body)
	}
}

func TestPageHandler_TokenExpired(t *testing.T) {
	req, err := http.NewRequest("GET", "/token-expired", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.PageHandler)
	handler.ServeHTTP(rr, req)

	// Check the response code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the title in the response body
	body := rr.Body.String()
	expectedTitle := "Token Expired - Trading Alerts"
	if !strings.Contains(body, `<title>`+expectedTitle+`</title>`) {
		t.Errorf("handler returned unexpected title: got %v", body)
	}
}

func TestPageHandler_Docs(t *testing.T) {
	req, err := http.NewRequest("GET", "/docs", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.PageHandler)
	handler.ServeHTTP(rr, req)

	// Check the response code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the title in the response body
	body := rr.Body.String()
	expectedTitle := "Documentation - Trading Alerts"
	if !strings.Contains(body, `<title>`+expectedTitle+`</title>`) {
		t.Errorf("handler returned unexpected title: got %v", body)
	}
}
