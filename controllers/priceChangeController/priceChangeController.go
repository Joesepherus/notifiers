package priceChangeController

import (
	"log"
	"net/http"
	"sort"
	"time"
	"tradingalerts/services/yahooService"
)

type stockDataHourly struct {
	high      float64
	low       float64
	timestamp int64
	change    float64
}

func changeHourToFilipino(hour int) int {
	filipino := hour + 12

	if filipino >= 24 {
		return filipino - 24
	}
	return filipino
}

func GetHourlyChange(w http.ResponseWriter, r *http.Request) {
	stockData, err := yahooService.GetStockCurrentValue(yahooService.YahooBaseURL, "USDJPY=X", "1h", "30d")
	log.Println("price: ", stockData, err)
	highs := stockData.Chart.Result[0].Indicators.Quote[0].High
	lows := stockData.Chart.Result[0].Indicators.Quote[0].Low
	timestamps := stockData.Chart.Result[0].Timestamp
	var hourlyData = []stockDataHourly{}

	for index, _ := range highs {
		hourlyData = append(hourlyData, stockDataHourly{high: highs[index], low: lows[index], timestamp: timestamps[index], change: highs[index] - lows[index]})
	}

	for _, data := range hourlyData {
		formattedTime := time.Unix(data.timestamp, 0).Format("2006-01-02 15:04:05")

		log.Printf("High: %.2f, Low: %.2f, Timestamp: %s, Change: %.2f\n", data.high, data.low, formattedTime, data.change)
	}

	var changeData = make(map[int][]float64)
	var changeDataAverage = make(map[int]float64)

	for _, data := range hourlyData {
		t := time.Unix(data.timestamp, 0)

		hour := t.Hour()
		changeData[hour] = append(changeData[hour], data.change)
	}

	for index, data := range changeData {
		log.Printf("Index: %d High: %d, \n", index, len(data))
		var average float64 = 0
		for _, dataInside := range data {
			average += dataInside
		}
		changeDataAverage[index] = average / float64(len(data))
	}

	keys := make([]int, 0, len(changeDataAverage))
	for key := range changeDataAverage {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	for index, key := range keys {
		value := changeDataAverage[key]
		filipinoHour := changeHourToFilipino(index)

		log.Printf("Index: %d, filipino: %d, Change: %.5f\n", key, filipinoHour, value) // Log the key as the index
	}

	// Basic health check response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Get Hourly Change"))
}
