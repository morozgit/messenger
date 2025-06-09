package handlers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

var (
	messages []Message
	mu       sync.Mutex
)

func GetMessages(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(http.StatusOK, messages)
}

func PostMessage(c *gin.Context) {
	var msg Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message"})
		return
	}

	mu.Lock()
	messages = append(messages, msg)
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"status": "Message received"})
}
