package main

import (
	outfit_recommender "AI_Agents/outfit-recommender"
	"AI_Agents/tools"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Load clothing rules into Milvus on startup
	err := outfit_recommender.LoadClothingRules()
	if err != nil {
		log.Printf("Failed to load clothing rules: %v", err)
	}

	// Pre-load IP database
	err = tools.InitLocationTool("outfit-recommender/data/IP2LOCATION-LITE-DB5_CN.CSV")
	if err != nil {
		log.Printf("Failed to initialize location tool: %v", err)
	}

	// http.HandleFunc("/ai_agents/outfit_recommend", outfitRecommendHandler)

	// port := "8081"
	// fmt.Printf("Starting server on port %s\n", port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
	outfit_recommender.StartServer()
}
