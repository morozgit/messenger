package handlers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Author    string `json:"author"`
	Content   string `json:"content"`
	Recipient string `json:"recipient"`
}

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan WSMessage)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clientsMu = sync.Mutex{}
)

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	for {
		var msg WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			clientsMu.Lock()
			delete(clients, conn)
			clientsMu.Unlock()
			break
		}

		// Отправляем сообщение пользователя всем
		broadcast <- msg

		// Если автор не бот, спрашиваем бота
		if msg.Author != "bot" {
			reply, err := askBot(msg.Content)
			if err != nil {
				reply = "Bot error: " + err.Error()
			}

			botMsg := WSMessage{
				Author:    "bot",
				Content:   reply,
				Recipient: msg.Author, // или по логике
			}
			broadcast <- botMsg
		}
	}
}

func StartBroadcast() {
	for {
		msg := <-broadcast
		clientsMu.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		clientsMu.Unlock()
	}
}
