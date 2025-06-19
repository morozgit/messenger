package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type BotRequest struct {
	Text string `json:"text"`
}

type BotResponse struct {
	Reply string `json:"reply"`
}

func askAiBot(text string) (string, error) {
	return askBotWithURL(text, os.Getenv("AI_BOT_CHAT_URL"))
}

func askMyAiBot(text string) (string, error) {
	return askBotWithURL(text, os.Getenv("MY_AI_BOT_CHAT_URL"))
}

func askBotWithURL(text, url string) (string, error) {
	reqBody := BotRequest{Text: text}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
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
