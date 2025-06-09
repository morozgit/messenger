// @title Messenger API
// @version 1.0
// @description Simple messaging API
// @host localhost:8080
// @BasePath /
package main

import (
	"messenger-backend/db"
	"messenger-backend/handlers"

	_ "messenger-backend/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	db.InitDB()
	go handlers.StartBroadcast()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/messenger/api/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/login", handlers.Login)
	r.GET("/ws", handlers.HandleWebSocket)
	r.GET("/users", handlers.GetUsers)
	r.POST("/add_users", handlers.AddUser)

	r.Run(":8080")
}
