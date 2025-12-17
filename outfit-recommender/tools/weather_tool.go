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
	Daily []struct {
		TempMax string `json:"tempMax"`
		TempMin string `json:"tempMin"`
		TextDay string `json:"textDay"`
	} `json:"daily"`
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

func GetWeather(location string) (avgTemp, maxTemp, minTemp float64, weather string, err error) {
	err = godotenv.Load("outfit-recommender/.env")
	if err != nil {
		return 0, 0, 0, "", fmt.Errorf("error loading .env file: %v", err)
	}

	apiKey := os.Getenv("WEATHER_API_TOKEN")
	if apiKey == "" {
		return 0, 0, 0, "", fmt.Errorf("WEATHER_API_TOKEN not set")
	}

	cityID, exists := cityIDs[location]
	if !exists {
		cityID = location // fallback to direct use, assuming it's already an ID
	}

	url := fmt.Sprintf("https://api.qweather.com/v7/weather/3d?location=%s&key=%s", cityID, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, 0, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, 0, 0, "", fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, 0, "", err
	}

	var weatherResp WeatherResponse
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		return 0, 0, 0, "", err
	}

	if weatherResp.Code != "200" {
		return 0, 0, 0, "", fmt.Errorf("API request failed with code %s", weatherResp.Code)
	}

	if len(weatherResp.Daily) == 0 {
		return 0, 0, 0, "", fmt.Errorf("no daily weather data")
	}

	tempMax, err := strconv.ParseFloat(weatherResp.Daily[0].TempMax, 64)
	if err != nil {
		return 0, 0, 0, "", err
	}
	tempMin, err := strconv.ParseFloat(weatherResp.Daily[0].TempMin, 64)
	if err != nil {
		return 0, 0, 0, "", err
	}
	avgTemp = (tempMax + tempMin) / 2
	maxTemp = tempMax
	minTemp = tempMin
	weather = weatherResp.Daily[0].TextDay

	return avgTemp, maxTemp, minTemp, weather, nil
}
