package main

import (
	outfit_recommender "AI_Agents/outfit_recommender"
	"AI_Agents/tools"
	"log"

	"github.com/gin-gonic/gin"
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
	err = tools.InitLocationTool("outfit_recommender/data/IP2LOCATION-LITE-DB5_CN.CSV")
	if err != nil {
		log.Printf("Failed to initialize location tool: %v", err)
	}
	StartServer()
}

func StartServer() {
	router := gin.Default()
	router.Use(RateLimit())
	router.GET("/ai_agents/outfit_recommend", outfit_recommender.OutfitRecommendHandler)
	router.Run(":8081")
}
