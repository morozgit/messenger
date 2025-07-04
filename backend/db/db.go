package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	UserPool *pgxpool.Pool
	ChatPool *pgxpool.Pool
)

func InitUserDB(connectURL string) {
	var err error
	UserPool, err = pgxpool.New(context.Background(), connectURL)
	if err != nil {
		log.Fatalf("Unable to connect to user database: %v\n", err)
	}

	if err = UserPool.Ping(context.Background()); err != nil {
		log.Fatalf("User database ping failed: %v\n", err)
	}

	_, err = UserPool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);
	`)
	if err != nil {
		log.Fatalf("Failed to initialize user schema: %v\n", err)
	}
}

func InitChatDB(connectURL string) {
	var err error
	ChatPool, err = pgxpool.New(context.Background(), connectURL)
	if err != nil {
		log.Fatalf("Unable to connect to chat database: %v\n", err)
	}

	if err = ChatPool.Ping(context.Background()); err != nil {
		log.Fatalf("Chat database ping failed: %v\n", err)
	}

	_, err = ChatPool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		chat_name TEXT NOT NULL,
		sender_username TEXT NOT NULL,
		content TEXT NOT NULL,
		sent_at TIMESTAMP WITH TIME ZONE DEFAULT now()
	);
	`)
	if err != nil {
		log.Fatalf("Failed to initialize chat schema: %v\n", err)
	}
}
