package yahooService

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"tradingalerts/types/yahooTypes"
)

const YahooBaseURL = "https://query1.finance.yahoo.com"

func GetStockCurrentValue(baseURL, symbol, interval, rangeData string) (*yahooTypes.StockResponse, error) {
	yahooFinanceUrl := fmt.Sprintf("%s/v8/finance/chart/%s?region=US&lang=en-US&includePrePost=false&interval=%s&useYfid=true&range=%s&corsDomain=finance.yahoo.com&.tsrc=finance", baseURL, symbol, interval, rangeData)

	resp, err := http.Get(yahooFinanceUrl)
	if err != nil {
		return nil, fmt.Errorf("error fetching stock price: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var stockData yahooTypes.StockResponse
	err = json.Unmarshal(body, &stockData)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	// Check if there is an error in the API response
	if stockData.Chart.Error.Code != "" {
		return nil, fmt.Errorf("yahoo API error: %s - %s", stockData.Chart.Error.Code, stockData.Chart.Error.Description)
	}

	return &stockData, nil
}
