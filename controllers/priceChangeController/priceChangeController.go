package priceChangeController

import (
	"log"
	"math"
	"net/http"
	"sort"
	"time"
	"tradingalerts/services/yahooService"
	"tradingalerts/templates"
	"tradingalerts/utils/errorUtils"
)

type stockDataHourly struct {
	high             float64
	low              float64
	timestamp        int64
	change           float64
	percentageChange float64
}

func changeHourToFilipino(hour int) int {
	filipino := hour + 12

	if filipino >= 24 {
		return filipino - 24
	}
	return filipino
}

func GetHourlyChange(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)
	// Extract form values
	symbol := r.FormValue("symbol")
	interval := r.FormValue("interval")

	stockData, err := yahooService.GetStockCurrentValue(yahooService.YahooBaseURL, symbol, "1h", interval)
	log.Println("price: ", stockData, err)
	highs := stockData.Chart.Result[0].Indicators.Quote[0].High
	lows := stockData.Chart.Result[0].Indicators.Quote[0].Low
	timestamps := stockData.Chart.Result[0].Timestamp
	var hourlyData = []stockDataHourly{}

	for index, _ := range highs {
		change := highs[index] - lows[index]
		percentageChange := (change / lows[index]) * 100
		hourlyData = append(hourlyData, stockDataHourly{
			high:             highs[index],
			low:              lows[index],
			timestamp:        timestamps[index],
			change:           change,
			percentageChange: percentageChange,
		})
	}

	var changeData = make(map[int][]float64)
	var changeDataAverage = make(map[int]float64)

	var changeDataPercentage = make(map[int][]float64)
	var changeDataPercentageAverage = make(map[int]float64)

	for _, data := range hourlyData {
		formattedTime := time.Unix(data.timestamp, 0).Format("2006-01-02 15:04:05")

		log.Printf("High: %.2f, Low: %.2f, Timestamp: %s, Change: %.2f\n", data.high, data.low, formattedTime, data.change)
		t := time.Unix(data.timestamp, 0)

		hour := t.Hour()
		changeData[hour] = append(changeData[hour], data.change)
		if !math.IsNaN(data.percentageChange) {
			changeDataPercentage[hour] = append(changeDataPercentage[hour], data.percentageChange)
		}
	}

	for index, data := range changeData {
		log.Printf("Index: %d High: %d, \n", index, len(data))
		var average float64 = 0
		for _, dataInside := range data {
			average += dataInside
		}
		changeDataAverage[index] = average / float64(len(data))
	}

	for hour, percentageChanges := range changeDataPercentage {
		var averagePercentage float64
		for _, percentage := range percentageChanges {
			averagePercentage += percentage
		}
		changeDataPercentageAverage[hour] = averagePercentage / float64(len(percentageChanges))
	}

	keys := make([]int, 0, len(changeDataAverage))
	for key := range changeDataAverage {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	for index, key := range keys {
		averageChange := changeDataAverage[key]
		averagePercentageChange := changeDataPercentageAverage[key]
		filipinoHour := changeHourToFilipino(index)

		log.Printf("Hour: %d, Filipino Hour: %d, Avg Change: %.5f, Avg Percentage Change: %.5f%%\n", key, filipinoHour, averageChange, averagePercentageChange)
	}

	formattedData := []struct {
		Hour                    int
		FilipinoHour            int
		ChangeAverage           float64
		PercentageChangeAverage float64
	}{
		// Store sorted and formatted change data
	}

	for _, hour := range keys {
		filipinoHour := changeHourToFilipino(hour) // Assuming this function converts to Filipino time
		formattedData = append(formattedData, struct {
			Hour                    int
			FilipinoHour            int
			ChangeAverage           float64
			PercentageChangeAverage float64
		}{
			Hour:                    hour,
			FilipinoHour:            filipinoHour,
			ChangeAverage:           changeDataAverage[hour],
			PercentageChangeAverage: changeDataPercentageAverage[hour],
		})
	}

	// Prepare data for the template
	templateLocation := templates.BaseLocation + "/price-change.html"
	data := map[string]interface{}{
		"Title":         "Price Change - Trading Alerts",
		"Content":       templateLocation,
		"Symbol":        symbol,
		"Interval":      interval,
		"FormattedData": formattedData,
	}

	templates.RenderTemplate(w, r, templateLocation, data)
}
