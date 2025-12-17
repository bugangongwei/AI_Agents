package outfit_recommender

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"outfit-recommender/outfit-recommender/tools"
)

/*
现在我想完善一下这个系统：我现在本地部署了一个Milvus StandAlone版本，这是前提；然后我需要细化这个系统：（1）用户自然语言输入的上下文信息比如“我不喜欢穿红色”，“除非温度低于5度，否则我不想在健身房穿长袖”等需要作为输出答案的上下文文参考，这就要求要有上下文的概念，以及上下文能够进行存储和过期（2）存储外部的时尚博客、专业建议、季节潮流等，作为 Agent 推荐的外部知识源，这里要求能够插入最少100条比较优质的数据到向量数据库里面，来模拟一个真实系统，如果模拟这些数据比较难，
*/

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
		ruleText := fmt.Sprintf("Temperature %d-%d°C, weather %s, preference %s: %s", rule.TemperatureMin, rule.TemperatureMax, rule.Weather, rule.Preference, rule.Outfit)
		err = tools.EmbedAndStoreRule(ruleText, rule.TemperatureMin, rule.TemperatureMax, rule.Weather, rule.Preference, rule.Outfit)
		if err != nil {
			log.Printf("Error storing rule: %v", err)
		}
	}
	return nil
}

func Start(args []string) (string, error) {
	// Set up flags with provided args
	oldArgs := os.Args
	os.Args = append([]string{"start"}, args...)
	defer func() { os.Args = oldArgs }()
	userInput := flag.String("question", "", "User question text")
	preference := flag.String("pref", "casual", "User preference (e.g., casual, formal)")
	location := flag.String("loc", "Beijing", "Location for weather")
	flag.Parse()

	if *userInput == "" {
		return "", fmt.Errorf("user question is required")
	}

	// Get weather
	avgTemp, _, _, weather, err := tools.GetWeather(*location)
	if err != nil {
		log.Printf("Error getting weather: %v, using defaults", err)
		avgTemp = 20
		weather = "sunny"
	}

	// Search for similar clothing rules in Milvus
	clothingRules, err := tools.SearchSimilar(*userInput, avgTemp, weather, *preference, 3)
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
	prompt := fmt.Sprintf("User question: %s\n\nRelevant clothing rules:\n%s\n\nProvide a personalized outfit recommendation based on the question and rules.", *userInput, rulesStr)

	// Get LLM recommendation
	recommendation, err := tools.GetLLMRecommendation(prompt)
	if err != nil {
		log.Printf("Error getting LLM recommendation: %v", err)
		recommendation = "Unable to generate recommendation"
	}

	return recommendation, nil
}

// func matchRule(rules []ClothingRule, temp int, weather, schedule, pref string) string {
// 	for _, rule := range rules {
// 		if temp >= rule.TemperatureMin && temp <= rule.TemperatureMax &&
// 			strings.ToLower(weather) == strings.ToLower(rule.Weather) &&
// 			strings.Contains(strings.ToLower(schedule), strings.ToLower(rule.Schedule)) &&
// 			strings.ToLower(pref) == strings.ToLower(rule.Preference) {
// 			return rule.Outfit
// 		}
// 	}
// 	return "T-shirt and jeans" // default
// }
