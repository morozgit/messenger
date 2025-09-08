package handlers

import (
	"context"
	"fmt"
	"messenger-backend/db"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

type ChatHistory struct {
	ID        int    `json:"id"`
	ChatID    int    `json:"chat_id"`
	SenderID  int    `json:"sender_id"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
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

func GetChatHistory(c *gin.Context) {
	chatIDStr := c.Query("chat_id")
	author := c.Query("author")
	recipient := c.Query("recipient")
	fmt.Printf("GetChatHistory called %v, %v, %v\n", chatIDStr, author, recipient)
	var chatID int
	var err error

	if chatIDStr != "" {
		chatID, err = strconv.Atoi(chatIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id"})
			return
		}
	} else if author != "" && recipient != "" {
		fmt.Println("else if")
		err = db.Pool.QueryRow(context.Background(),
			`SELECT chat_id FROM chats 
     			WHERE (user1 = (SELECT user_id FROM users WHERE username = $1) 
        		AND user2 = (SELECT user_id FROM users WHERE username = $2))
        		OR (user1 = (SELECT user_id FROM users WHERE username = $2) 
        		AND user2 = (SELECT user_id FROM users WHERE username = $1))
     			LIMIT 1`,
			author, recipient).Scan(&chatID)
		if err != nil {
			fmt.Println("else if if")
			c.JSON(http.StatusOK, []ChatHistory{})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need chat_id or author+recipient"})
		return
	}

	rows, err := db.Pool.Query(context.Background(),
		`SELECT m.id, m.chat_id, m.sender_id, u.username AS author, m.content, m.created_at
     FROM messages m
     JOIN users u ON u.user_id = m.sender_id
     WHERE m.chat_id = $1
     ORDER BY m.created_at ASC`, chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	fmt.Println("Resolved chatID:", chatID)
	var history []ChatHistory
	for rows.Next() {
		var createdAt time.Time
		var m ChatHistory
		if err := rows.Scan(&m.ID, &m.ChatID, &m.SenderID, &m.Content, &m.Author, &createdAt); err == nil {
			m.CreatedAt = createdAt.Format(time.RFC3339)
			history = append(history, m)
		}
	}
	fmt.Println("Chat history:", history)

	c.JSON(http.StatusOK, history)
}
