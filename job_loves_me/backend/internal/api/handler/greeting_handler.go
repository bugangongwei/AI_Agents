package handler

import (
	"net/http"

	"github.com/bugangongwei/job_loves_me/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type GreetingHandler struct {
	service *service.GreetingService
}

func NewGreetingHandler(resumeServ *service.ResumeService) *GreetingHandler {
	return &GreetingHandler{service: service.NewGreetingService(resumeServ)}
}

func (h *GreetingHandler) GenerateGreeting(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		JDText string `json:"jd_text" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For MVP, using sync version
	greeting, err := h.service.GenerateGreetingSync(c.Request.Context(), userID.(uint), req.JDText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate greeting: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"greeting": greeting,
	})
}
