package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// LLMClient LLM客户端
type LLMClient struct {
	APIKey      string
	Model       string
	Temperature float64
	MaxTokens   int
	BaseURL     string // 自定义API地址（支持DeepSeek、Ollama等）
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// NewLLMClient 创建LLM客户端
func NewLLMClient(apiKey, model string, temperature float64, maxTokens int, baseURL string) *LLMClient {
	return &LLMClient{
		APIKey:      apiKey,
		Model:       model,
		Temperature: temperature,
		MaxTokens:   maxTokens,
		BaseURL:     baseURL,
	}
}

// Chat 发送聊天请求
func (c *LLMClient) Chat(messages []ChatMessage) (string, error) {
	if c.APIKey == "" {
		return "", errors.New("API Key未配置")
	}

	// 检测是否使用Ollama
	if strings.Contains(c.BaseURL, "localhost:11434") || c.APIKey == "ollama" {
		return c.chatWithOllama(messages)
	}

	// 云端API请求
	request := ChatRequest{
		Model:       c.Model,
		Messages:    messages,
		Temperature: c.Temperature,
		MaxTokens:   c.MaxTokens,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("构造请求失败: %w", err)
	}

	// 使用自定义BaseURL或默认OpenAI URL
	apiURL := "https://api.openai.com/v1/chat/completions"
	if c.BaseURL != "" {
		// 支持DeepSeek、智谱AI、通义千问、Groq等
		if strings.HasSuffix(c.BaseURL, "/v1") || strings.HasSuffix(c.BaseURL, "/openai") {
			apiURL = c.BaseURL + "/chat/completions"
		} else {
			apiURL = c.BaseURL + "/v1/chat/completions"
		}
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API返回错误 %d: %s", resp.StatusCode, string(body))
	}

	var response ChatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", errors.New("API返回空响应")
	}

	return response.Choices[0].Message.Content, nil
}

// chatWithOllama Ollama专用请求
func (c *LLMClient) chatWithOllama(messages []ChatMessage) (string, error) {
	type OllamaRequest struct {
		Model    string                 `json:"model"`
		Messages []ChatMessage          `json:"messages"`
		Stream   bool                   `json:"stream"`
		Options  map[string]interface{} `json:"options,omitempty"`
	}

	request := OllamaRequest{
		Model:    c.Model,
		Messages: messages,
		Stream:   false,
		Options: map[string]interface{}{
			"temperature": c.Temperature,
			"num_predict": c.MaxTokens,
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("构造请求失败: %w", err)
	}

	apiURL := c.BaseURL + "/api/chat"

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 300 * time.Second} // Ollama可能较慢，增加超时时间
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求Ollama失败: %w (请确保Ollama服务正在运行)", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama返回错误 %d: %s", resp.StatusCode, string(body))
	}

	// Ollama响应格式
	type OllamaResponse struct {
		Model     string `json:"model"`
		CreatedAt string `json:"created_at"`
		Message   struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Done bool `json:"done"`
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("解析Ollama响应失败: %w", err)
	}

	return ollamaResp.Message.Content, nil
}

// CompletionWithJSON 使用JSON格式响应
func (c *LLMClient) CompletionWithJSON(systemPrompt, userPrompt string) (map[string]interface{}, error) {
	messages := []ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	response, err := c.Chat(messages)
	if err != nil {
		return nil, err
	}

	// 解析JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("解析JSON响应失败: %w", err)
	}

	return result, nil
}
