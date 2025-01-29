package webapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MessagePayload struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type MessageSenderWebAPI struct {
	URL string
}

func NewMessageSenderWebAPI(url string) *MessageSenderWebAPI {
	return &MessageSenderWebAPI{
		URL: url,
	}
}
func (wc *MessageSenderWebAPI) SendMessage(to, content string) (string, error) {
	payload := MessagePayload{
		To:      to,
		Content: content,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}
	req, err := http.NewRequest("POST", wc.URL, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("unexpected response code: %d, body: %s", resp.StatusCode, string(body))
	}

	var responseMap map[string]interface{}
	if err := json.Unmarshal(body, &responseMap); err != nil {
		return "", fmt.Errorf("failed to parse response body: %w", err)
	}

	messageID, ok := responseMap["messageId"].(string)
	if !ok {
		return "", fmt.Errorf("messageId not found in response")
	}

	return messageID, nil
}
