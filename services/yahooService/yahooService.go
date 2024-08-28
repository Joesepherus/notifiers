package yahooService

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"notifiers/types/yahooTypes"
)

func GetStockCurrentValue(symbol string) (*yahooTypes.StockResponse, error) {
	yahooFinanceUrl := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?region=US&lang=en-US&includePrePost=false&interval=2m&useYfid=true&range=1d&corsDomain=finance.yahoo.com&.tsrc=finance", symbol)

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

	return &stockData, nil
}
