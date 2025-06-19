package handlers

import (
	"fmt"
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
		fmt.Println("HandleWebSocket", msg)
		// Если автор не бот, спрашиваем бота
		if msg.Author != "Ai_Bot" && msg.Author != "My_Ai_Bot" {
			var reply string
			var err error

			switch msg.Recipient {
			case "Ai_Bot":
				reply, err = askAiBot(msg.Content)
			case "My_Ai_Bot":
				reply, err = askMyAiBot(msg.Content)
			default:
				reply = "Неизвестный получатель. Доступны: Ai_Bot, My_Ai_Bot"
			}

			if err != nil {
				reply = "Ошибка бота: " + err.Error()
			}

			botMsg := WSMessage{
				Author:    msg.Recipient,
				Content:   reply,
				Recipient: msg.Author,
			}
			broadcast <- botMsg
		}
	}
}

func StartBroadcast() {
	for {
		msg := <-broadcast
		fmt.Println("StartBroadcast", msg)
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
