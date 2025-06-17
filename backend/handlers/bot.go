package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type BotRequest struct {
	Text string `json:"text"`
}

type BotResponse struct {
	Reply string `json:"reply"`
}

func askBot(text string) (string, error) {
	reqBody := BotRequest{Text: text}
	jsonData, _ := json.Marshal(reqBody)

	_ = godotenv.Load()
	botAPIURL := os.Getenv("botAPIURL")
	if botAPIURL == "" {
		log.Fatal("botAPIURL is not set")
	}

	resp, err := http.Post(botAPIURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var botResp BotResponse
	err = json.NewDecoder(resp.Body).Decode(&botResp)
	if err != nil {
		return "", err
	}
	return botResp.Reply, nil
}
