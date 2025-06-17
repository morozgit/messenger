package main

import (
	"context"
	"log"
	"messenger-backend/db"
	"messenger-backend/handlers"

	_ "messenger-backend/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ensureBotUser() error {
	_, err := db.Pool.Exec(context.Background(),
		"INSERT INTO users(username, password) VALUES($1, $2) ON CONFLICT (username) DO NOTHING",
		"Bot", "some_secure_password")
	return err
}

func main() {
	db.InitDB()
	if err := ensureBotUser(); err != nil {
		log.Fatalf("failed to ensure bot user: %v", err)
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
	}

	r.Run(":8080")
}
