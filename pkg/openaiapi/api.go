package openaiapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CompletionResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
	Model string `json:"model"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Index        int         `json:"index"`
		Message      ChatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func (c *Client) GetCompletion(prompt string, maxTokens int, temperature float32, model string) (string, error) {
	apiURL := "https://api.openai.com/v1/completions"
	data := map[string]interface{}{
		"model":       model,
		"prompt":      prompt,
		"max_tokens":  maxTokens,
		"temperature": temperature,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("API request failed")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var completionResponse CompletionResponse
	err = json.Unmarshal(body, &completionResponse)
	if err != nil {
		return "", err
	}

	return completionResponse.Choices[0].Text, nil
}

type EmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Model string `json:"model"`
}

func (c *Client) GetEmbedding(paragraph string) ([]float32, error) {
	apiURL := "https://api.openai.com/v1/embeddings" // Update the API endpoint
	data := map[string]string{
		"input": paragraph,
		"model": "text-embedding-ada-002",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("API request failed")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var embeddingResponse EmbeddingResponse
	err = json.Unmarshal(body, &embeddingResponse)
	if err != nil {
		return nil, err
	}

	return embeddingResponse.Data[0].Embedding, nil
}

//Its is develop
func (c *Client) GetChatCompletion(messages string, temperature float32, model string) (string, error) {
	apiURL := "https://api.openai.com/v1/chat/completions"
	data := map[string]interface{}{
		"model": model,
		"messages": ChatMessage{
			Role:    "user",
			Content: messages,
		},
		// "temperature": temperature,

	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("API request failed")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var chatCompletionResponse ChatCompletionResponse
	err = json.Unmarshal(body, &chatCompletionResponse)
	if err != nil {
		return "", err
	}

	return chatCompletionResponse.Choices[0].Message.Content, nil
}
