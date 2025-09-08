package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB(connectURL string) {
	var err error
	Pool, err = pgxpool.New(context.Background(), connectURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	if err = Pool.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
	}

	// users
	_, err = Pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
	);
	`)
	if err != nil {
		log.Fatalf("Failed to initialize users schema: %v\n", err)
	}

	// chats
	_, err = Pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS chats (
		chat_id SERIAL PRIMARY KEY,
		user1 INTEGER NOT NULL REFERENCES users(user_id),
		user2 INTEGER NOT NULL REFERENCES users(user_id),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
	);
	`)
	if err != nil {
		log.Fatalf("Failed to initialize chats schema: %v\n", err)
	}

	// messages
	_, err = Pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		chat_id INTEGER NOT NULL REFERENCES chats(chat_id),
		sender_id INTEGER NOT NULL REFERENCES users(user_id),
		content TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		is_read BOOLEAN DEFAULT FALSE
	);
	`)
	if err != nil {
		log.Fatalf("Failed to initialize messages schema: %v\n", err)
	}
}
