package handlers

import (
	"context"
	"log"
	"messenger-backend/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 400,401 {object} map[string]string
// @Router /login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		log.Println("BindJSON error or empty fields:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing username or password"})
		return
	}

	var user User
	err := db.UserPool.QueryRow(context.Background(),
		`SELECT id, username FROM users WHERE username = $1 AND password = $2`,
		req.Username, req.Password).
		Scan(&user.Id, &user.Username)

	if err != nil {
		log.Println("Login failed:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Welcome, " + user.Username,
	})
}

// AddUser godoc
// @Summary Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User data"
// @Success 201
// @Failure 400,500 {object} map[string]string
// @Router /add_users [post]
func AddUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.BindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing username or password"})
		return
	}

	_, err := db.UserPool.Exec(context.Background(),
		"INSERT INTO users(username, password) VALUES($1, $2) ON CONFLICT DO NOTHING",
		req.Username, req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// GetUsers godoc
// @Summary Get list of users
// @Tags users
// @Produce json
// @Success 200 {array} User
// @Failure 500 {object} map[string]string
// @Router /users [get]
func GetUsers(c *gin.Context) {
	rows, err := db.UserPool.Query(context.Background(), "SELECT id, username, password FROM users ORDER BY username")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Username, &u.Password); err == nil {
			users = append(users, u)
		}
	}

	c.JSON(http.StatusOK, users)
}
