package yahooService

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStockCurrentValue_Success(t *testing.T) {
	symbol := "AAPL"
	stockData, err := GetStockCurrentValue(YahooBaseURL, symbol, "2m", "1d")

	retunedSymbol := stockData.Chart.Result[0].Meta.Symbol
	assert.NoError(t, err)
	assert.Equal(t, symbol, retunedSymbol)
}

func TestGetStockCurrentValue_Fail(t *testing.T) {
	_, err := GetStockCurrentValue(YahooBaseURL, "AAAPL", "2m", "1d")

	assert.EqualError(t, err, "yahoo API error: Not Found - No data found, symbol may be delisted")
}

func TestMockGetStockCurrentValue_Success(t *testing.T) {
	// Create a new test server that returns a predefined response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		mockResponse := `{
			"chart": {
				"result": [{
					"meta": {
						"symbol": "AAPL",
						"regularMarketPrice": 150.00
					}
				}]
			}
		}`
		w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	// Call the function you want to test with the mock server URL
	stockData, err := GetStockCurrentValue(ts.URL, "AAPL", "2m", "1d")

	// Check the results
	assert.NoError(t, err)
	assert.Equal(t, "AAPL", stockData.Chart.Result[0].Meta.Symbol)
	assert.Equal(t, 150.00, stockData.Chart.Result[0].Meta.RegularMarketPrice)
}

func TestMockGetStockCurrentValue_Fail(t *testing.T) {
	// Create a new test server that returns a predefined error response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		mockResponse := `{
			"chart": {
				"result": null,
				"error": {
					"code": "Not Found",
					"description": "No data found, symbol may be delisted"
				}
			}
		}`
		w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	// Call the function you want to test with the mock server URL
	_, err := GetStockCurrentValue(ts.URL, "AAAPL", "2m", "1d")

	// Check the results
	assert.EqualError(t, err, "yahoo API error: Not Found - No data found, symbol may be delisted")
}
