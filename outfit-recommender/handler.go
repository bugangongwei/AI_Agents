package outfit_recommender

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func outfitRecommendHandler(c *gin.Context) {
	question := c.Query("question")
	pref := c.Query("pref")
	location := c.Query("loc")

	if question == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'question' parameter"})
		return
	}

	if pref == "" {
		pref = "casual"
	}
	if location == "" {
		location = "Shanghai"
	}

	recommendation, err := GetRecommendation(question, pref, location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error: %v", err)})
		return
	}

	response := map[string]string{"recommendation": recommendation}
	c.JSON(http.StatusOK, response)
}

func StartServer() {
	router := gin.Default()
	router.GET("/ai_agents/outfit_recommend", outfitRecommendHandler)
	router.Run(":8081")
}
