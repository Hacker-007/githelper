package commit

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLMResponse struct {
	Message LLMMessage `json:"message"`
}

type LLMClient struct {
	messageHistory []LLMMessage
}

func (client *LLMClient) SendStandaloneMessage(prompt string) (string, error) {
	options, buffer := map[string]interface{}{
		"model":  "llama3.1",
		"prompt": prompt,
		"stream": false,
	}, new(bytes.Buffer)

	if err := json.NewEncoder(buffer).Encode(options); err != nil {
		return "", err
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", buffer)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	var llmResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&llmResp); err != nil {
		return "", err
	}

	return llmResp["response"].(string), nil
}

func (client *LLMClient) SendChatMessage(message LLMMessage) (*LLMMessage, error) {
	client.messageHistory = append(client.messageHistory, message)
	options, buffer := map[string]interface{}{
		"model":    "llama3.1",
		"messages": client.messageHistory,
		"stream":   false,
	}, new(bytes.Buffer)

	if err := json.NewEncoder(buffer).Encode(options); err != nil {
		return nil, err
	}

	resp, err := http.Post("http://localhost:11434/api/chat", "application/json", buffer)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var llmResp LLMResponse
	if err := json.NewDecoder(resp.Body).Decode(&llmResp); err != nil {
		return nil, err
	}

	return &llmResp.Message, nil
}
