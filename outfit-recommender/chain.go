package outfit_recommender

import (
	"AI_Agents/outfit-recommender/tools"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
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
	rulesData, err := os.ReadFile("outfit-recommender/data/clothing_rules.json")
	if err != nil {
		return fmt.Errorf("failed to load rules: %v", err)
	}
	var rules Rules
	err = json.Unmarshal(rulesData, &rules)
	if err != nil {
		return fmt.Errorf("failed to unmarshal rules: %v", err)
	}

	// Store rules in vector database
	for _, rule := range rules.Rules {
		log.Printf("Storing rule: %+v", rule)
		ruleText := fmt.Sprintf("温度 %d-%d°C, 天气 %s, 偏好 %s: %s", rule.TemperatureMin, rule.TemperatureMax, rule.Weather, rule.Preference, rule.Outfit)
		err = tools.EmbedAndStoreRule(ruleText, rule.TemperatureMin, rule.TemperatureMax, rule.Weather, rule.Preference, rule.Outfit)
		if err != nil {
			log.Printf("Error storing rule: %v", err)
		}
	}
	return nil
}

// GetRecommendation generates outfit recommendations based on user input, preferences, and location.
// It integrates weather data, searches for relevant clothing rules, and uses LLM for personalized suggestions.
func GetRecommendation(userInput, preference, location string) (string, error) {
	if userInput == "" {
		return "", fmt.Errorf("user input is required")
	}

	// Retrieve current weather conditions
	maxTemp, minTemp, weather, err := tools.GetWeather(location)
	if err != nil {
		log.Printf("Failed to retrieve weather data: %v. Using default values (min: 15°C, max: 25°C, weather: 晴天)", err)
		minTemp, maxTemp = 15, 25 // Corrected: min should be lower than max
		weather = "晴天"
	}

	// Search for similar clothing rules in the vector database
	clothingRules, err := tools.SearchSimilar(userInput, maxTemp, minTemp, weather, preference, 3)
	if err != nil {
		log.Printf("Error searching for similar clothing rules: %v. Proceeding with empty rules", err)
		clothingRules = []string{} // Fallback to empty slice
	}

	// Efficiently build the rules string using strings.Builder for better performance
	var rulesBuilder strings.Builder
	for _, rule := range clothingRules {
		rulesBuilder.WriteString(rule)
		rulesBuilder.WriteString("\n")
	}
	rulesStr := rulesBuilder.String()

	// Construct the prompt for the LLM, including weather context and instructions for natural language output with reasoning
	prompt := fmt.Sprintf("用户问题: %s\n\n当前天气: 温度 %d-%d°C, 天气 %s\n\n相关穿衣规则:\n%s\n\n请基于用户问题、当前天气和相关穿衣规则，提供个性化的穿衣推荐。请解释您的决策过程，包括参考的天气条件、用户要求和穿衣规则。推荐应使用自然语言，避免特殊符号，直接面向用户。", userInput, minTemp, maxTemp, weather, rulesStr)

	// Obtain recommendation from LLM
	recommendation, err := tools.GetLLMRecommendation(prompt)
	if err != nil {
		log.Printf("Failed to get LLM recommendation: %v", err)
		// Return a meaningful error instead of a generic string, allowing caller to handle appropriately
		return "", fmt.Errorf("unable to generate recommendation: %w", err)
	}

	return recommendation, nil
}

func ClearOldData() error {
	return tools.ClearCollection()
}
