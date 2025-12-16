package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type WeatherResponse struct {
	Code string `json:"code"`
	Now  struct {
		Temp string `json:"temp"`
		Text string `json:"text"`
	} `json:"now"`
}

var cityIDs = map[string]string{
	"Beijing":   "101010100",
	"Shanghai":  "101020100",
	"Guangzhou": "101280101",
	"Shenzhen":  "101280601",
	"Hangzhou":  "101210101",
	"Nanjing":   "101190101",
	"Wuhan":     "101200101",
	"Chengdu":   "101270101",
	"Chongqing": "101040100",
	"Xi'an":     "101110101",
}

func GetWeather(location string) (temperature float64, weather string, err error) {
	err = godotenv.Load("outfit-recommender/.env")
	if err != nil {
		return 0, "", fmt.Errorf("error loading .env file: %v", err)
	}

	apiKey := os.Getenv("WEATHER_API_TOKEN")
	if apiKey == "" {
		return 0, "", fmt.Errorf("WEATHER_API_TOKEN not set")
	}

	cityID, exists := cityIDs[location]
	if !exists {
		cityID = location // fallback to direct use, assuming it's already an ID
	}

	url := fmt.Sprintf("https://api.qweather.com/v7/weather/now?location=%s&key=%s", cityID, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, "", fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}

	var weatherResp WeatherResponse
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		return 0, "", err
	}

	if weatherResp.Code != "200" {
		return 0, "", fmt.Errorf("API request failed with code %s", weatherResp.Code)
	}

	tempFloat, err := strconv.ParseFloat(weatherResp.Now.Temp, 64)
	if err != nil {
		return 0, "", err
	}
	temperature = tempFloat
	weather = weatherResp.Now.Text

	return temperature, weather, nil
}
