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
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/messenger/api/", func(c *gin.Context) {
		c.Redirect(302, "/messenger/api/index.html")
	})
	r.GET("/messenger/api/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/messenger/api/swagger/doc.json"),
	))

	r.POST("/login", handlers.Login)
	r.GET("/ws", handlers.HandleWebSocket)
	r.GET("/users", handlers.GetUsers)
	r.POST("/add_users", handlers.AddUser)

	r.Run(":8080")
}
