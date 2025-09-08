package handlers

import (
	"context"
	"fmt"
	"log"
	"messenger-backend/db"
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

type BroadcastMessage struct {
	ChatID   int    `json:"chat_id"`
	SenderID int    `json:"sender_id"`
	Content  string `json:"content"`
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

	fmt.Println("Клиент подключился")
	defer fmt.Println("Клиент отключился")

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

		broadcast <- msg

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

		// Получаем sender_id по имени автора
		var senderID int
		err := db.Pool.QueryRow(context.Background(),
			"SELECT user_id FROM users WHERE username = $1", msg.Author).Scan(&senderID)
		if err != nil {
			log.Println("Не найден sender_id для:", msg.Author, err)
			continue
		}

		// Получаем chat_id по участникам (или создаём, если нет)
		var chatID int
		err = db.Pool.QueryRow(context.Background(),
			`SELECT chat_id FROM chats 
             WHERE (user1 = (SELECT user_id FROM users WHERE username = $1) 
                AND user2 = (SELECT user_id FROM users WHERE username = $2))
                OR (user1 = (SELECT user_id FROM users WHERE username = $2) 
                AND user2 = (SELECT user_id FROM users WHERE username = $1))
             LIMIT 1`,
			msg.Author, msg.Recipient).Scan(&chatID)
		if err != nil {
			// Если чата нет, создаём
			var user1, user2 int
			_ = db.Pool.QueryRow(context.Background(), "SELECT user_id FROM users WHERE username = $1", msg.Author).Scan(&user1)
			_ = db.Pool.QueryRow(context.Background(), "SELECT user_id FROM users WHERE username = $1", msg.Recipient).Scan(&user2)
			err = db.Pool.QueryRow(context.Background(),
				"INSERT INTO chats(user1, user2) VALUES($1, $2) RETURNING chat_id", user1, user2).Scan(&chatID)
			if err != nil {
				log.Println("Ошибка создания чата:", err)
				continue
			}
		}

		// Вставляем сообщение
		_, err = db.Pool.Exec(
			context.Background(),
			`INSERT INTO messages (chat_id, sender_id, content) VALUES ($1, $2, $3)`,
			chatID, senderID, msg.Content,
		)
		if err != nil {
			log.Println("Ошибка при вставке сообщения:", err)
		}

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
