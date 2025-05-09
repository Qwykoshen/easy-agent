package controllers

import (
	"easy-agent/config"
	"easy-agent/services/deepseek"
	"net/http"

	"github.com/gin-gonic/gin"
)

var deepseekClient *deepseek.Client

// InitDeepSeekClient initializes the DeepSeek API client
func InitDeepSeekClient() {
	deepseekClient = deepseek.NewClient(config.DeepSeekAPIKey)
}

// ChatRequest represents the request body for chat API
type ChatRequest struct {
	Messages    []deepseek.ChatMessage `json:"messages" binding:"required"`
	Model       string                 `json:"model" binding:"required"`
	Temperature float64                `json:"temperature,omitempty"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
}

// Chat handles the chat completion request
func Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.Messages) >= 3 {
		req.Messages = req.Messages[len(req.Messages)-3:]
	}
	chatReq := &deepseek.ChatRequest{
		Model:       req.Model,
		Messages:    req.Messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
	}
	if chatReq.Model == "" {
		chatReq.Model = "deepseek-chat"
	}

	resp, err := deepseekClient.CreateChatCompletion(chatReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
