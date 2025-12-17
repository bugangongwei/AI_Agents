package outfit_recommender

import (
	"AI_Agents/outfit-recommender/tools"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type ClothingRule struct {
	TemperatureMin int    `json:"temperature_min"`
	TemperatureMax int    `json:"temperature_max"`
	Weather        string `json:"weather"`
	Schedule       string `json:"schedule"`
	Preference     string `json:"preference"`
	Outfit         string `json:"outfit"`
}

type Rules struct {
	Rules []ClothingRule `json:"rules"`
}

func LoadClothingRules() error {
	rulesData, err := ioutil.ReadFile("outfit-recommender/data/clothing_rules.json")
	if err != nil {
		return fmt.Errorf("failed to load rules: %v", err)
	}
	var rules Rules
	json.Unmarshal(rulesData, &rules)

	// Store rules in vector database
	for _, rule := range rules.Rules {
		log.Printf("Storing rule: %+v", rule)
		ruleText := fmt.Sprintf("Temperature %d-%d°C, weather %s, preference %s: %s", rule.TemperatureMin, rule.TemperatureMax, rule.Weather, rule.Preference, rule.Outfit)
		err = tools.EmbedAndStoreRule(ruleText, rule.TemperatureMin, rule.TemperatureMax, rule.Weather, rule.Preference, rule.Outfit)
		if err != nil {
			log.Printf("Error storing rule: %v", err)
		}
	}
	return nil
}

func GetRecommendation(userInput, preference, location string) (string, error) {
	if userInput == "" {
		return "", fmt.Errorf("user question is required")
	}

	// Get weather
	maxTemp, minTemp, weather, err := tools.GetWeather(location)
	if err != nil {
		log.Printf("Error getting weather: %v, using defaults(15, 25, sunny)", err)
		maxTemp = 15
		minTemp = 25
		weather = "sunny"
	}

	// Search for similar clothing rules in Milvus
	clothingRules, err := tools.SearchSimilar(userInput, maxTemp, minTemp, weather, preference, 3)
	if err != nil {
		log.Printf("Error searching similar rules: %v", err)
		clothingRules = []string{}
	}

	// Combine rules into a string
	rulesStr := ""
	for _, rule := range clothingRules {
		rulesStr += rule + "\n"
	}

	// Build prompt with question and rules
	prompt := fmt.Sprintf("User question: %s\n\nRelevant clothing rules:\n%s\n\nProvide a personalized outfit recommendation based on the question and rules.", userInput, rulesStr)

	// Get LLM recommendation
	recommendation, err := tools.GetLLMRecommendation(prompt)
	if err != nil {
		log.Printf("Error getting LLM recommendation: %v", err)
		recommendation = "Unable to generate recommendation"
	}

	return recommendation, nil
}
