package handler

import (
	"net/http"

	"github.com/bugangongwei/job_loves_me/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ResumeHandler struct {
	service *service.ResumeService
}

func NewResumeHandler() *ResumeHandler {
	return &ResumeHandler{service: service.NewResumeService()}
}

func (h *ResumeHandler) UploadResume(c *gin.Context) {
	userID, _ := c.Get("user_id")
	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No resume file uploaded"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer f.Close()

	resume, err := h.service.SaveAndParseResume(userID.(uint), file.Filename, f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save and parse resume: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Resume uploaded and parsed successfully",
		"resume":  resume,
	})
}

func (h *ResumeHandler) GetLatestResume(c *gin.Context) {
	userID, _ := c.Get("user_id")
	resume, err := h.service.GetLatestResume(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resume not found"})
		return
	}

	c.JSON(http.StatusOK, resume)
}
