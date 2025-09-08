package main

import (
	"context"
	"log"
	"messenger-backend/db"
	"messenger-backend/handlers"
	"os"

	_ "messenger-backend/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ensureAiBotUser() error {
	_, err := db.Pool.Exec(context.Background(),
		"INSERT INTO users(username, password) VALUES($1, $2) ON CONFLICT (username) DO NOTHING",
		"Ai_Bot", "some_secure_password")
	return err
}
func ensureMyAiBotUser() error {
	_, err := db.Pool.Exec(context.Background(),
		"INSERT INTO users(username, password) VALUES($1, $2) ON CONFLICT (username) DO NOTHING",
		"My_Ai_Bot", "some_secure_password")
	return err
}

func main() {
	_ = godotenv.Load()
	url_db := os.Getenv("DATABASE_URL")
	if url_db == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db.InitDB(url_db)
	if err := ensureAiBotUser(); err != nil {
		log.Fatalf("failed to ensure Ai_Bot user: %v", err)
	}
	if err := ensureMyAiBotUser(); err != nil {
		log.Fatalf("failed to ensure My_Ai_Bot user: %v", err)
	}
	go handlers.StartBroadcast()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/messenger/docs/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/messenger/docs/swagger/doc.json"),
	))

	api := r.Group("/messenger/api")
	{
		api.POST("/login", handlers.Login)
		api.POST("/add_users", handlers.AddUser)
		api.GET("/users", handlers.GetUsers)
		api.GET("/ws", handlers.HandleWebSocket)
		api.GET("/chat_history", handlers.GetChatHistory)
	}

	r.Run(":8080")
}
