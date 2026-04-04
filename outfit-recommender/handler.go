package outfit_recommender

import (
	"AI_Agents/tools"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func outfitRecommendHandler(c *gin.Context) {
	question := c.Query("question")
	pref := c.Query("pref")
	location := c.Query("loc")
	schedule := c.Query("schedule")

	if question == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'question' parameter"})
		return
	}

	if pref == "" {
		pref = "casual"
	}
	if location == "" {
		// If location is missing, resolve city from IP
		ip := c.ClientIP()
		location = tools.GetCityFromIP(ip)
		if location == "" {
			location = "Shanghai"
		}
	}

	recommendation, err := GetRecommendation(question, pref, location, schedule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error: %v", err)})
		return
	}

	response := map[string]string{"recommendation": recommendation}
	c.JSON(http.StatusOK, response)
}

func StartServer() {
	router := gin.Default()
	router.Use(RateLimit())
	router.GET("/ai_agents/outfit_recommend", outfitRecommendHandler)
	router.Run(":8081")
}
