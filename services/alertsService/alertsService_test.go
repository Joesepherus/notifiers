package alertsService

import (
	"errors"
	"testing"

	// Adjust this import path as needed
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAlerts_Success(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	SetDB(db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock the expected query and the returned rows
	rows := sqlmock.NewRows([]string{"id", "symbol", "trigger_value", "alert_type"}).
		AddRow(1, "AAPL", 150.00, "lower").
		AddRow(2, "GOOGL", 2800.00, "higher")

	mock.ExpectQuery("SELECT id, symbol, trigger_value, alert_type FROM alerts WHERE triggered = FALSE ORDER BY symbol").
		WillReturnRows(rows)

	// Call the function you're testing
	alerts, _ := GetAlerts()
	// // Check if the function works as expected
	assert.NoError(t, err)
	assert.Len(t, alerts, 2)
	assert.Equal(t, "AAPL", alerts[0].Symbol)
	assert.Equal(t, "GOOGL", alerts[1].Symbol)
}

func TestGetAlerts_QueryError(t *testing.T) {
	// Create mock DB and mock query results
	db, mock, err := sqlmock.New()
	SetDB(db)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock the expected query to return an error
	mock.ExpectQuery("SELECT id, symbol, trigger_value, alert_type FROM alerts WHERE triggered = FALSE ORDER BY symbol").
		WillReturnError(errors.New("query error"))

	// Call the function you're testing
	alerts, err := GetAlerts()

	// Check if the error is handled correctly
	assert.Nil(t, alerts)
	assert.EqualError(t, err, "failed to query alerts: query error")
}

func TestGetAlerts_ScanError(t *testing.T) {
	// Create mock DB and mock query results
	db, mock, err := sqlmock.New()
	SetDB(db)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock the expected query and a faulty row (mismatching columns or types)
	rows := sqlmock.NewRows([]string{"id", "symbol", "trigger_value", "alert_type"}).
		AddRow("invalid", "AAPL", 150.00, "lower")

	mock.ExpectQuery("SELECT id, symbol, trigger_value, alert_type FROM alerts WHERE triggered = FALSE ORDER BY symbol").
		WillReturnRows(rows)

	// Call the function you're testing
	alerts, err := GetAlerts()

	// Check if the scanning error is handled correctly
	assert.Nil(t, alerts)
	assert.Contains(t, err.Error(), "failed to scan row")
}
